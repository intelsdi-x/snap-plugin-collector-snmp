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

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/configReader"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/serror"
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
)

type Plugin struct {
	metricConfig configReader.Metrics
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
	)
}

// New creates initialized instance of snmp collector
func New() *Plugin {
	return &Plugin{metricConfig: make(configReader.Metrics, 0)}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (p *Plugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	item, err := config.GetConfigItem(cfg, setFileConfigVar)
	if err != nil {
		return nil, err
	}
	setFilePath, ok := item.(string)
	if !ok {
		return nil, serror.New(fmt.Errorf("Incorrect type of configuration variable, cannot parse value of %s to string", setFileConfigVar), nil)
	}

	_, serr := configReader.GetMetricsConfig(setFilePath)
	if serr != nil {
		log.WithFields(serr.Fields()).Error(serr.Error())
		return nil, serr
	}

	mts := []plugin.MetricType{}
	//TODO: building namespaces from configuration

	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (p *Plugin) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {

	config, _ := getConfigItems(metricTypes[0])

	_, serr := configReader.GetSnmpAgentConfig(config)
	if serr != nil {
		return nil, serr
	}

	mts := []plugin.MetricType{}
	//TODO: collection of metrics

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

//getConfigItems reads agent parameters from configuration and creates map of agent parameters
func getConfigItems(metricTypes plugin.MetricType) (map[string]interface{}, error) {
	cfg := make(map[string]interface{})
	for _, snmpAgentParam := range configReader.SnmpAgentConfigParameters {
		item, err := config.GetConfigItem(metricTypes, snmpAgentParam)
		if err != nil {
			serr := serror.New(fmt.Errorf("Missing parameter in configuration (%s)", snmpAgentParam))
			log.WithFields(serr.Fields()).Warn(serr.Error())
		} else {
			cfg[snmpAgentParam] = item
		}
	}
	return cfg, nil
}
