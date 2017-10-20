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

	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/configReader"
	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/snmp"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/core"
	"github.com/k-sone/snmpgo"
	log "github.com/sirupsen/logrus"
)

const (
	// PluginName namespace part
	PluginName = "snmp"

	// Version of plugin
	Version = 2

	// Vendor namespace part
	Vendor = "intel"

	// setFileConfigVar configuration variable to define path to setfile
	setFileConfigVar = "setfile"

	// tagSnmpAgentName indicates SNMP agent name, tag which is added to metrics
	tagSnmpAgentName = "SNMP_AGENT_NAME"

	// tagSnmpAgentAddress indicates SNMP agent address, tag which is added to metrics
	tagSnmpAgentAddress = "SNMP_AGENT_ADDRESS"

	// tagOid indicates metric OID, tag which is added to metrics
	tagOid = "OID"

	// the max time a connection can be unused.
	connectionIdle = time.Minute * 30

	// the time between checking of connections usage
	connectionWait = time.Minute * 15
)

// Plugin main structure
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
	newHandler(hostConfig configReader.SnmpAgent) (*snmpgo.SNMP, error)
	readElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, error)
}

var (
	snmp_              = snmpInterface(&snmpType{})
	snmpConnections    = make(map[string]connection)
	mtxSnmpConnections = &sync.Mutex{}
)

func init() {
	go watchConnections()
}

// New creates initialized instance of snmp collector
func New() *Plugin {
	return &Plugin{metricsConfigs: make(map[string]configReader.Metric)}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *Plugin) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	configs, err := getMetricsConfig(cfg)
	if err != nil {
		return nil, err
	}

	mts := []plugin.Metric{}
	for _, cfg := range configs {

		namespace := plugin.NewNamespace(Vendor, PluginName)
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

			mt := plugin.Metric{
				Namespace:   namespace,
				Description: cfg.Description,
				Unit:        cfg.Unit,
			}
			mts = append(mts, mt)

		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *Plugin) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	var mtxMetrics sync.Mutex
	var wgCollectedMetrics sync.WaitGroup

	//initialization of plugin structure (only once)
	if !p.initialized {
		configs, err := getMetricsConfig(metrics[0].Config)
		if err != nil {
			return nil, err
		}
		for _, cfg := range configs {
			namespace := core.NewNamespace(Vendor, PluginName)
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

	agentConfig, err := configReader.GetSnmpAgentConfig(metrics[0].Config)
	if err != nil {
		return nil, err
	}

	//lock using of connections in watchConnections
	mtxSnmpConnections.Lock()
	defer mtxSnmpConnections.Unlock()

	conn, err := getConnection(agentConfig)
	if err != nil {
		return nil, err
	}

	mts := []plugin.Metric{}

	for _, metric := range metrics {

		//get metrics to collect
		metricsConfigs, err := getMetricsToCollect(metric.Namespace.String(), p.metricsConfigs)
		if err != nil {
			return nil, err
		}

		wgCollectedMetrics.Add(len(metricsConfigs))

		for _, cfg := range metricsConfigs {

			go func(cfg configReader.Metric) {

				defer wgCollectedMetrics.Done()

				conn.mtx.Lock()

				//get value of metric/metrics
				results, err := snmp_.readElements(conn.handler, cfg.Oid, cfg.Mode)
				if err != nil {
					log.Warn(err)
					conn.mtx.Unlock()
					return
				}

				//get dynamic elements of namespace parts
				err = getDynamicNamespaceElements(conn.handler, results, &cfg)
				if err != nil {
					conn.mtx.Unlock()
					return
				}

				conn.lastUsed = time.Now()
				conn.mtx.Unlock()

				for i, result := range results {

					//build namespace for metric
					namespace := plugin.NewNamespace(Vendor, PluginName)
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
					val, err := convertSnmpDataToMetric(result.Variable.String(), result.Variable.Type())
					if err != nil {
						continue
					}

					//modify numeric metric - use scale and shift parameters
					data := modifyNumericMetric(val, cfg.Scale, cfg.Shift)

					//creating metric
					mt := plugin.Metric{
						Namespace: namespace,
						Data:      data,
						Timestamp: time.Now(),
						Tags: map[string]string{
							tagSnmpAgentName:    agentConfig.Name,
							tagSnmpAgentAddress: agentConfig.Address,
							tagOid:              result.Oid.String()},
						Unit:        metric.Unit,
						Description: metric.Description,
					}

					//adding metric to list of metrics
					mtxMetrics.Lock()

					//filter specific instance
					nsPattern := strings.Replace(metric.Namespace.String(), "*", ".*", -1)
					matched, err := regexp.MatchString(nsPattern, mt.Namespace.String())
					if err != nil {
						logFields := map[string]interface{}{"namespace": mt.Namespace.String(), "pattern": nsPattern, "match_error": err}
						err := fmt.Errorf("Cannot parse namespace element for matching")
						log.WithFields(logFields).Warn(err)
						return
					}
					if matched {
						mts = append(mts, mt)
					}

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
func (p *Plugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	err := policy.AddNewStringRule([]string{Vendor, PluginName}, setFileConfigVar, true)
	if err != nil {
		return *policy, err
	}

	return *policy, nil
}

//NewHandler creates new connection with SNMP agent
func (s *snmpType) newHandler(hostConfig configReader.SnmpAgent) (*snmpgo.SNMP, error) {
	return snmp.NewHandler(hostConfig)
}

//ReadElements reads data using SNMP requests
func (s *snmpType) readElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, error) {
	return snmp.ReadElements(handler, oid, mode)
}

//getConnection gets connection with SNMP agent, checks if connection with specified SNMP agent exists, if not a new connection is initialized
func getConnection(agentConfig configReader.SnmpAgent) (connection, error) {
	if conn, ok := snmpConnections[agentConfig.Address]; ok {
		return conn, nil
	}
	handler, err := snmp_.newHandler(agentConfig)
	if err != nil {
		return connection{}, err
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
func getDynamicNamespaceElements(handler *snmpgo.SNMP, results []*snmpgo.VarBind, metric *configReader.Metric) error {
	for i := 0; i < len(metric.Namespace); i++ {
		//clear slice with dynamic parts of namespace
		metric.Namespace[i].Values = []string{}

		switch metric.Namespace[i].Source {

		case configReader.NsSourceString:
			continue

		case configReader.NsSourceSNMP:
			parts, err := snmp_.readElements(handler, metric.Namespace[i].Oid, metric.Mode)
			if err != nil {
				return err
			}
			for _, part := range parts {
				metricNamePart := ns.ReplaceNotAllowedCharsInNamespacePart(part.Variable.String())
				metric.Namespace[i].Values = append(metric.Namespace[i].Values, metricNamePart)
			}

		case configReader.NsSourceIndex:
			for _, r := range results {
				oidParts := strings.Split(strings.Trim(r.Oid.String(), "."), ".")

				if uint(len(oidParts)) <= metric.Namespace[i].OidPart {

					logFields := log.Fields{
						"namespace_part_configuration": metric.Namespace[i],
						"oid_part":                     metric.Namespace[i].OidPart,
						"number_of_oid_elements":       len(metric.Namespace[i].Values)}
					err := fmt.Errorf("Incorrect value of `oid_part`  in configuration of namespace")
					log.WithFields(logFields).Warn(err)
					return err
				}
				metric.Namespace[i].Values = append(metric.Namespace[i].Values, oidParts[metric.Namespace[i].OidPart])
			}
		}

		if len(metric.Namespace[i].Values) != len(results) {
			logFields := log.Fields{
				"namespace_part_configuration": metric.Namespace[i],
				"number_of_results":            len(results),
				"number_of_namespace_elements": len(metric.Namespace[i].Values)}
			err := fmt.Errorf("Incorrect configuration of dynamic elements of namespace, number of namespace elements is not equal to number of results")
			log.WithFields(logFields).Warn(err)
			return err
		}
	}
	return nil
}

//getMetricsConfig reads metrics parameters from configuration
func getMetricsConfig(cfg plugin.Config) (configReader.Metrics, error) {
	setFilePath, err := cfg.GetString(setFileConfigVar)
	if err != nil {
		return nil, err
	}

	configs, err := configReader.GetMetricsConfig(setFilePath)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

//getMetricsToCollects gets configuration of metrics which are requested through task
func getMetricsToCollect(namespace string, metrics map[string]configReader.Metric) (map[string]configReader.Metric, error) {
	collectedMetrics := make(map[string]configReader.Metric)

	// Filter out setfile based metrics by given namespace
	for ns := range metrics {
		matched, err := regexp.MatchString(strings.Replace(ns, "*", ".*", -1), namespace)
		if err != nil {
			return nil, err
		}

		if matched {
			collectedMetrics[ns] = metrics[ns]
		}
	}
	if len(collectedMetrics) == 0 {
		return nil, fmt.Errorf("Metric namespace (`%s`) is not supported by this plugin", namespace)
	}
	return collectedMetrics, nil
}

//convertSnmpDataToMetric converts data received using SNMP request to supported data type
func convertSnmpDataToMetric(snmpData string, snmpType string) (interface{}, error) {
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
		log.WithFields(log.Fields{"data": snmpData, "type": snmpType}).Warn(
			fmt.Errorf("Unrecognized type of data, metric is returned as string"))
		val = snmpData
	}

	if err != nil {
		log.WithFields(log.Fields{"data": snmpData, "type": snmpType}).Warn(err)
		return nil, err
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
