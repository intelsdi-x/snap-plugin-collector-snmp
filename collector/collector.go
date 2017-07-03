// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package collector

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/configReader"
	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/snmp"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/serror"
	"github.com/k-sone/snmpgo"
)

const (
	//pluginName namespace part
	pluginName = "snmp"

	// version of plugin
	version = 1

	//pluginType type of plugin
	pluginType = plugin.CollectorPluginType

	//vendor namespace part
	vendor = "intel"

	//setFileConfigVar configuration variable to define path to setfile
	setFileConfigVar = "setfile"

	//tag_snmp_agent_name indicates SNMP agent name, tag which is added to metrics
	tag_snmp_agent_name = "SNMP_AGENT_NAME"

	//tag_snmp_agent_address indicates SNMP agent address, tag which is added to metrics
	tag_snmp_agent_address = "SNMP_AGENT_ADDRESS"

	//tag_oid indicates metric OID, tag which is added to metrics
	tag_oid = "OID"

	//the max time a connection can be unused.
	connectionIdle = time.Minute * 30

	//the time between checking of connections usage
	connectionWait = time.Minute * 15
)

type Plugin struct {
	initialized    bool
	metricsConfigs map[string]configReader.Metric
}

type connection struct {
	handler  *snmpgo.SNMP
	mtx      *sync.Mutex
	lastUsed time.Time
}

type snmpType struct{}

type snmpInterface interface {
	newHandler(hostConfig configReader.SnmpAgent) (*snmpgo.SNMP, serror.SnapError)
	readElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, serror.SnapError)
}

var (
	snmp_              = snmpInterface(&snmpType{})
	snmpConnections    = make(map[string]connection)
	mtxSnmpConnections = &sync.Mutex{}
)

func init() {
	go watchConnections()
}

// Meta returns plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginName,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.RoutingStrategy(plugin.StickyRouting),
		plugin.ConcurrencyCount(1),
	)
}

// New creates initialized instance of snmp collector
func New() *Plugin {
	return &Plugin{metricsConfigs: make(map[string]configReader.Metric)}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *Plugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	configs, serr := getMetricsConfig(cfg)
	if serr != nil {
		return nil, fmt.Errorf(serr.Error())
	}

	mts := []plugin.MetricType{}
	for _, cfg := range configs {

		namespace := core.NewNamespace(vendor, pluginName)
		for _, ns := range cfg.Namespace {
			if ns.Source == configReader.NsSourceString {
				namespace = namespace.AddStaticElement(ns.String)
			} else {
				namespace = namespace.AddDynamicElement(ns.Name, ns.Description)
			}
		}

		if _, metricExist := p.metricsConfigs[namespace.String()]; metricExist {
			logFields := map[string]interface{}{
				"namespace":                     namespace.String(),
				"previous_metric_configuration": p.metricsConfigs[namespace.String()],
				"current_metric_configuration":  cfg,
			}
			log.WithFields(logFields).Warn(fmt.Errorf("Plugin configuration file (`setfile`) contains metrics definitions which expose the same namespace, only one of them is in use. Correction of plugin configuration file (`setfile`) is recommended."))
		} else {
			p.metricsConfigs[namespace.String()] = cfg

			mt := plugin.MetricType{
				Namespace_:   namespace,
				Description_: cfg.Description,
				Unit_:        cfg.Unit,
			}
			mts = append(mts, mt)

		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *Plugin) CollectMetrics(metrics []plugin.MetricType) ([]plugin.MetricType, error) {
	var mtxMetrics sync.Mutex
	var wgCollectedMetrics sync.WaitGroup

	//initialization of plugin structure (only once)
	if !p.initialized {

		configs, serr := getMetricsConfig(metrics[0])
		if serr != nil {
			return nil, fmt.Errorf(serr.Error())
		}
		for _, cfg := range configs {
			namespace := core.NewNamespace(vendor, pluginName)
			for _, ns := range cfg.Namespace {
				if ns.Source == configReader.NsSourceString {
					namespace = namespace.AddStaticElement(ns.String)
				} else {
					namespace = namespace.AddDynamicElement(ns.Name, ns.Description)
				}
			}

			if _, metricExist := p.metricsConfigs[namespace.String()]; metricExist {
				logFields := map[string]interface{}{
					"namespace":                     namespace.String(),
					"previous_metric_configuration": p.metricsConfigs[namespace.String()],
					"current_metric_configuration":  cfg,
				}
				log.WithFields(logFields).Warn(fmt.Errorf("Plugin configuration file (`setfile`) contains metrics definitions which expose the same namespace, only one of them is in use. Correction of plugin configuration file (`setfile`) is recommended."))
			} else {
				//add metric configuration to plugin metric map
				p.metricsConfigs[namespace.String()] = cfg
			}
		}
		p.initialized = true
	}

	agentConfigVars := getAgentConfig(metrics[0])
	agentConfig, serr := configReader.GetSnmpAgentConfig(agentConfigVars)
	if serr != nil {
		return nil, serr
	}

	//lock using of connections in watchConnections
	mtxSnmpConnections.Lock()
	defer mtxSnmpConnections.Unlock()

	conn, serr := getConnection(agentConfig)
	if serr != nil {
		return nil, serr
	}

	mts := []plugin.MetricType{}

	for _, metric := range metrics {

		//get metrics to collect
		metricsConfigs, serr := getMetricsToCollect(metric.Namespace().String(), p.metricsConfigs)
		if serr != nil {
			log.WithFields(serr.Fields()).Warn(serr.Error())
			return nil, serr
		}

		wgCollectedMetrics.Add(len(metricsConfigs))

		for _, cfg := range metricsConfigs {

			go func(cfg configReader.Metric) {

				defer wgCollectedMetrics.Done()

				conn.mtx.Lock()

				//get value of metric/metrics
				results, serr := snmp_.readElements(conn.handler, cfg.Oid, cfg.Mode)
				if serr != nil {
					log.WithFields(serr.Fields()).Warn(serr.Error())
					conn.mtx.Unlock()
					return
				}

				resultMap, serr := snmpResults2Map(results, cfg.Oid)
				if serr != nil {
					log.WithFields(serr.Fields()).Warn(serr.Error())
					conn.mtx.Unlock()
					return
				}

				//get dynamic elements of namespace parts
				serr = getDynamicNamespaceElements(conn.handler, resultMap, &cfg)
				if serr != nil {
					log.WithFields(serr.Fields()).Warn(serr.Error())
					conn.mtx.Unlock()
					return
				}

				conn.lastUsed = time.Now()
				conn.mtx.Unlock()

				for i, result := range resultMap {

					//build namespace for metric
					namespace := core.NewNamespace(vendor, pluginName)
					offset := len(namespace)
					for j, ns := range cfg.Namespace {
						if ns.Source == configReader.NsSourceString {
							namespace = namespace.AddStaticElements(ns.String)
						} else {
							namespace = namespace.AddDynamicElement(ns.Name, ns.Description)
							namespace[j+offset].Value = ns.Values[i]
						}
					}

					//convert metric types
					val, serr := convertSnmpDataToMetric(result.Variable.String(), result.Variable.Type())
					if serr != nil {
						log.WithFields(serr.Fields()).Warn(serr.Error())
						continue
					}

					//modify numeric metric - use scale and shift parameters
					data := modifyNumericMetric(val, cfg.Scale, cfg.Shift)

					//creating metric
					mt := plugin.MetricType{
						Namespace_: namespace,
						Data_:      data,
						Timestamp_: time.Now(),
						Tags_: map[string]string{
							tag_snmp_agent_name:    agentConfig.Name,
							tag_snmp_agent_address: agentConfig.Address,
							tag_oid:                result.Oid.String()},
						Unit_:        metric.Unit(),
						Description_: metric.Description(),
					}

					//adding metric to list of metrics
					mtxMetrics.Lock()
					mts = append(mts, mt)
					mtxMetrics.Unlock()
				}
			}(cfg)
		}
		wgCollectedMetrics.Wait()
	}
	return mts, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (p *Plugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	rule, err := cpolicy.NewStringRule(setFileConfigVar, true)
	if err != nil {
		return cp, err
	}
	rule.Description = "Configuration file"
	config.Add(rule)

	return cp, nil
}

//NewHandler creates new connection with SNMP agent
func (s *snmpType) newHandler(hostConfig configReader.SnmpAgent) (*snmpgo.SNMP, serror.SnapError) {
	return snmp.NewHandler(hostConfig)
}

//ReadElements reads data using SNMP requests
func (s *snmpType) readElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, serror.SnapError) {
	return snmp.ReadElements(handler, oid, mode)
}

//getConnection gets connection with SNMP agent, checks if connection with specified SNMP agent exists, if not a new connection is initialized
func getConnection(agentConfig configReader.SnmpAgent) (connection, serror.SnapError) {
	if conn, ok := snmpConnections[agentConfig.Address]; ok {
		return conn, nil
	}
	handler, serr := snmp_.newHandler(agentConfig)
	if serr != nil {
		return connection{}, serr
	}
	snmpConnections[agentConfig.Address] = connection{handler: handler, mtx: &sync.Mutex{}}
	return snmpConnections[agentConfig.Address], nil
}

//watchConnections observes SNMP connections and closes unused connections
func watchConnections() {
	for {
		time.Sleep(connectionWait)
		mtxSnmpConnections.Lock()
		for k := range snmpConnections {
			if time.Now().Sub(snmpConnections[k].lastUsed) > connectionIdle {

				//close the connection
				snmpConnections[k].handler.Close()

				//remove the connection
				delete(snmpConnections, k)
			}
		}
		mtxSnmpConnections.Unlock()
	}
}

//getDynamicNamespaceElements gets dynamic elements of namespace, either sending SNMP requests or using part of OID
func getDynamicNamespaceElements(handler *snmpgo.SNMP, resultMap map[int]*snmpgo.VarBind, metric *configReader.Metric) serror.SnapError {
	for i := 0; i < len(metric.Namespace); i++ {
		//clear slice with dynamic parts of namespace
		metric.Namespace[i].Values = make(map[int]string)

		switch metric.Namespace[i].Source {

		case configReader.NsSourceString:
			continue

		case configReader.NsSourceSNMP:
			parts, serr := snmp_.readElements(handler, metric.Namespace[i].Oid, metric.Mode)

			if serr != nil {
				return serr
			}

			partsMap, serr := snmpResults2Map(parts, metric.Namespace[i].Oid)

			if serr != nil {
				return serr
			}

			for idx, part := range partsMap {
				if _, ok := resultMap[idx]; ok {
					metricNamePart := ns.ReplaceNotAllowedCharsInNamespacePart(part.Variable.String())
					metric.Namespace[i].Values[idx] = metricNamePart
				}
			}

		case configReader.NsSourceIndex:
			for _, r := range resultMap {
				oidParts := strings.Split(strings.Trim(r.Oid.String(), "."), ".")

				if uint(len(oidParts)) <= metric.Namespace[i].OidPart {

					logFields := map[string]interface{}{
						"namespace_part_configuration": metric.Namespace[i],
						"oid_part":                     metric.Namespace[i].OidPart,
						"number_of_oid_elements":       len(metric.Namespace[i].Values)}
					return serror.New(fmt.Errorf("Incorrect value of `oid_part`  in configuration of namespace"), logFields)
				}
				idx := oidParts[metric.Namespace[i].OidPart]
				iidx, _ := strconv.Atoi(idx) //TODO: error handling?
				metric.Namespace[i].Values[iidx] = idx
			}
		}
	}
	return nil
}

//getMetricsConfig reads metrics parameters from configuration
func getMetricsConfig(cfg interface{}) (configReader.Metrics, serror.SnapError) {
	item, err := config.GetConfigItem(cfg, setFileConfigVar)
	if err != nil {
		return nil, serror.New(err)
	}
	setFilePath, ok := item.(string)
	if !ok {
		return nil, serror.New(fmt.Errorf("Incorrect type of configuration variable, cannot parse value of %s to string", setFileConfigVar), nil)
	}

	configs, serr := configReader.GetMetricsConfig(setFilePath)
	if serr != nil {
		log.WithFields(serr.Fields()).Error(serr.Error())
		return nil, serr
	}

	return configs, nil
}

//getAgentConfig reads agent parameters from configuration and creates map of agent parameters
func getAgentConfig(metricTypes plugin.MetricType) map[string]interface{} {
	cfg := make(map[string]interface{})
	for _, snmpAgentParam := range configReader.SnmpAgentConfigParameters {
		item, err := config.GetConfigItem(metricTypes, snmpAgentParam)
		if err == nil {
			cfg[snmpAgentParam] = item
		}
	}
	return cfg
}

//getMetricsToCollects gets configuration of metrics which are requested through task
func getMetricsToCollect(namespace string, metrics map[string]configReader.Metric) (map[string]configReader.Metric, serror.SnapError) {
	collectedMetrics := make(map[string]configReader.Metric)

	// change `*` into regexp `.*` which matches any characters
	namespace = strings.Replace(namespace, "*", ".*", -1)

	for ns, _ := range metrics {
		matched, err := regexp.MatchString(namespace, ns)
		if err != nil {
			return nil, serror.New(err)
		}

		if matched {
			collectedMetrics[ns] = metrics[ns]
		}
	}
	if len(collectedMetrics) == 0 {
		return nil, serror.New(fmt.Errorf("Metric namespace (`%s`) is not supported by this plugin", namespace), nil)
	}
	return collectedMetrics, nil
}

//convertSnmpDataToMetric converts data received using SNMP request to supported data type
func convertSnmpDataToMetric(snmpData string, snmpType string) (interface{}, serror.SnapError) {
	var val interface{}
	var err error

	switch snmpType {
	case "Counter", "Counter32", "Gauge32", "UInteger32", "TimeTicks":
		val, err = strconv.ParseUint(snmpData, 10, 32)
	case "Counter64":
		val, err = strconv.ParseUint(snmpData, 10, 64)
	case "Integer", "Integer32":
		val, err = strconv.ParseInt(snmpData, 10, 32)
	case "OctetString", "IpAddress", "Object Identifier":
		val = snmpData
	default:
		serr := serror.New(fmt.Errorf("Unrecognized type of data, metric is returned as string"), map[string]interface{}{"data": snmpData, "type": snmpType})
		log.WithFields(serr.Fields()).Warn(serr.Error())
		val = snmpData
	}

	if err != nil {
		return nil, serror.New(err, map[string]interface{}{"data": snmpData, "type": snmpType})
	}
	return val, nil
}

//modifyNumericMetric modifies value of numeric data using scale and shift parameters
func modifyNumericMetric(data interface{}, scale float64, shift float64) interface{} {
	//check if scale and shift parameters are set
	if scale == 1.0 && shift == 0.0 {
		return data
	}
	var modifiedData interface{}
	//use shift and scale and convert data to float64
	switch data.(type) {
	case uint32:
		modifiedData = float64(data.(uint32))*scale + shift
	case uint64:
		modifiedData = float64(data.(uint64))*scale + shift
	case int:
		modifiedData = float64(data.(int))*scale + shift
	case int32:
		modifiedData = float64(data.(int32))*scale + shift
	case int64:
		modifiedData = float64(data.(int64))*scale + shift
	default:
		modifiedData = data
	}
	return modifiedData
}

func snmpResults2Map(results []*snmpgo.VarBind, oid string) (map[int]*snmpgo.VarBind, serror.SnapError) {
	index := ""
	res := make(map[int]*snmpgo.VarBind)
	oidTrimmed := strings.Trim(oid, ".")
	oidLen := len(strings.Split(oidTrimmed, "."))

	for _, r := range results {
		roidTrimmed := strings.Trim(r.Oid.String(), ".")
		roidLen := len(strings.Split(roidTrimmed, "."))
		if ((roidLen == oidLen) || (roidLen == oidLen+1)) && strings.HasPrefix(roidTrimmed, oidTrimmed) {
			if roidLen == oidLen+1 {
				index = strings.Split(roidTrimmed, ".")[oidLen]
			} else {
				index = strings.Split(roidTrimmed, ".")[oidLen-1]
			}
			idx, _ := strconv.Atoi(index) //TODO: error handling?
			res[idx] = r
		} else {
			return nil, serror.New(fmt.Errorf("Inconsistent Oid in response"), nil)
		}
	}
	return res, nil
}
