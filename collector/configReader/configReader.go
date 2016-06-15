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

package configReader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/intelsdi-x/snap/core/serror"
	"github.com/mitchellh/mapstructure"
)

const (
	//ModeSingle option in mode of metric
	ModeSingle = "single"

	//ModeMultiple option in mode of metric
	ModeMultiple = "multiple"

	//prefixSuffixSource indicates source of prefix or suffix in configuration of metric
	PrefixSuffixSource = "source"

	//SourceSNMP option in source of prefix or suffix
	SourceSNMP = "snmp"

	//SourceString option in source of prefix or suffix
	SourceString = "string"

	//PrefixSuffixOID indicates OID which is in use to get prefix or suffix for metric name
	PrefixSuffixOID = "OID"

	//PrefixSuffixString indicates value which is addes as prefix or suffix for metric name
	PrefixSuffixString = "string"

	//hostName indicates SNMP host name
	hostName = "snmp_host_name"

	//hostAddress indicates SNMP host address
	hostAddress = "snmp_host_address"

	//hostSnmpVersion indicates SNMP version in host configuration
	hostSnmpVersion = "snmp_version"

	//hostCommunity indicates community (SNMP  v1 &  SNMP v2c) in host configuration
	hostCommunity = "community"

	//hostNetwork indicates network which is used in host configuration, see net.Dial parameter
	hostNetwork = "network"

	//hostUserName indicates user name (SNMP v3) in host configuration
	hostUserName = "user_name"

	//hostSecurityLevel indicates security level (SNMP v3) in host configuration
	hostSecurityLevel = "security_level"

	//hostAuthPassword indicates authentication protocol pass phrase (SNMP v3) in host configuration
	hostAuthPassword = "auth_protocol"

	//hostAuthProtocol indicates authentication protocol (SNMP v3) in host configuration
	hostAuthProtocol = "auth_protocol"

	//hostPrivPassword indicates privacy protocol pass phrase (SNMP v3) in host configuration
	hostPrivPassword = "priv_password"

	//hostPrivProtocol indicates privacy protocol (SNMP v3) in host configuration
	hostPrivProtocol = "priv_protocol"

	//hostSecurityEngineId indicates security engine ID (SNMP v3) in host configuration
	hostSecurityEngineId = "security_engine_id"

	//hostContextEngineID indicates context engine ID (SNMP v3) in host configuration
	hostContextEngineID = "context_engine_id"

	//hostContextName indicates context name (SNMP v3) in host configuration
	hostContextName = "context_name"

	//hostRetries indicates number of connection retries in host configuration
	hostRetries = "retries"

	//hostTimeout indicates timeout for network connection in host configuration
	hostTimeout = "timeout"

	//metricOid indicates OID which is use to receive metric
	metricOid = "OID"

	//metricMode indicates mode of metric
	metricMode = "mode"

	//metricPrefix indicates prefix configuration in metric configuration
	metricPrefix = "prefix"

	//metircSuffix indicates suffix configuration in metric configuration
	metricSuffix = "suffix"

	//metricScale indicates scale value which can be used to multiplication of metric value
	metricScale = "scale"

	//snmpv1 name of SNMP v1 in configuration
	snmpv1 = "v1"

	//snmpv2 name of SNMP v2c in configuration
	snmpv2 = "v2c"

	//snmpv3 symbol of SNMP v3 in configuration
	snmpv3 = "v3"

	//defaultRetries default number of connection retries
	defaultRetries = 1

	//defaultTimeout timeout for network connection
	defaultTimeout = 5

	//missingRequiredParameter error message for missing required parameter
	missingRequiredParameter = "Missing required parameter in host configuration"

	//inCorrectValueOfParameter error message for incorrect value of parameter
	inCorrectValueOfParameter = "Incorrect value of parameter (%s), possible options: %v"
)

type Host struct {
	Name             string `mapstructure:"snmp_host_name"`
	SnmpVersion      string `mapstructure:"snmp_version"`
	Address          string `mapstructure:"snmp_host_address"`
	Community        string `mapstructure:"community"`
	Network          string `mapstructure:"network"`
	UserName         string `mapstructure:"user_name"`
	SecurityLevel    string `mapstructure:"security_level"`
	AuthPassword     string `mapstructure:"auth_password"`
	AuthProtocol     string `mapstructure:"auth_protocol"`
	PrivPassword     string `mapstructure:"priv_password"`
	PrivProtocol     string `mapstructure:"priv_protocol"`
	SecurityEngineId string `mapstructure:"security_engine_id"`
	ContextEngineId  string `mapstructure:"context_engine_id"`
	ContextName      string `mapstructure:"context_name"`
	Retries          uint   `mapstructure:"retries"`
	Timeout          int    `mapstructure:"timeout"`
}

type Metric struct {
	Oid         string            `json:"OID"`
	Mode        string            `json:"mode"`
	Prefix      map[string]string `json:"prefix"`
	Suffix      map[string]string `json:"suffix"`
	Unit        string            `json:"unit"`
	Description string            `json:"description"`
	Shift       float64           `json:"shift"`
	Scale       float64           `json:"scale"`
}

type Metrics map[string]Metric

type cfgReaderType struct{}

type reader interface {
	ReadFile(s string) ([]byte, error)
}

var (
	//HostConfigParameters slice of host configuration parameters
	HostConfigParameters = []string{hostName, hostAddress, hostSnmpVersion, hostCommunity, hostNetwork,
		hostUserName, hostSecurityLevel, hostAuthPassword, hostAuthProtocol, hostPrivPassword,
		hostPrivProtocol, hostSecurityEngineId, hostContextEngineID, hostContextName, hostRetries, hostTimeout}

	//modeOptions slice of options for mode parameter
	modeOptions = []interface{}{ModeSingle, ModeMultiple}

	//sourcePrefixSuffixOptions slice of options for source of prefix or suffix in configuration
	sourcePrefixSuffixOptions = []interface{}{SourceSNMP, SourceString}

	//snmpVersionOptions slice of options for SNMP version
	snmpVersionOptions = []interface{}{snmpv1, snmpv2, snmpv3}

	//communityOptions slice of options for SNMP community
	communityOptions = []interface{}{"public", "private"}

	//securityLevelOptions slice of options for SNMP security level
	securityLevelOptions = []interface{}{"NoAuthNoPriv", "AuthNoPriv", "AuthPriv"}

	//authProtocolOptions slice of options for SNMP authentication protocol
	authProtocolOptions = []interface{}{"MD5", "SHA"}

	//privProtocolOptions slice of options for SNMP privacy protocol
	privProtocolOptions = []interface{}{"DES", "AES"}

	//cfgReader provides possibility to read metric configuration from file or from different source
	cfgReader = reader(&cfgReaderType{})
)

func (r *cfgReaderType) ReadFile(s string) ([]byte, error) {
	return ioutil.ReadFile(s)
}

//GetMetricsConfig decodes and validates configuration of SNMP host
func GetHostConfig(configMap map[string]interface{}) (Host, serror.SnapError) {
	config, serr := decodeHostConfig(configMap)
	if serr != nil {
		return config, serr
	}

	serr = validateHostConfig(config)
	if serr != nil {
		return config, serr
	}

	return config, nil
}

//GetMetricsConfig reads and validates configuration of metrics
func GetMetricsConfig(setFilePath string) (Metrics, serror.SnapError) {
	config, serr := readMetricConfigFile(setFilePath)
	if serr != nil {
		return config, serr
	}

	serr = validateMetricConfig(config)
	if serr != nil {
		return config, serr
	}

	return config, nil
}

//decodeHostConfig decodes configuration of SNMP host into structure
func decodeHostConfig(config map[string]interface{}) (Host, serror.SnapError) {
	var hostConfig Host
	logFields := map[string]interface{}{}
	err := mapstructure.Decode(config, &hostConfig)
	if err != nil {
		return hostConfig, serror.New(err, logFields)
	}
	return hostConfig, nil
}

//validateMetricConfig validates configuration of SNMP host
func validateHostConfig(config Host) serror.SnapError {
	logFields := map[string]interface{}{}
	logFields["host_config"] = config

	if !checkSetParameter(config.Address) {
		logFields["parameter"] = hostAddress
		return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
	}

	if !checkSetParameter(config.SnmpVersion) {
		logFields["parameter"] = hostSnmpVersion
		return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
	}

	if !checkPossibleOptions(config.SnmpVersion, snmpVersionOptions) {
		logFields["parameter"] = hostSnmpVersion
		return serror.New(fmt.Errorf(inCorrectValueOfParameter, config.SnmpVersion, snmpVersionOptions), logFields)
	}

	if config.SnmpVersion == snmpv1 || config.SnmpVersion == snmpv2 {
		//check required fields for SNMP v1 and SNMP v2c
		if !checkSetParameter(config.SnmpVersion) {
			logFields["parameter"] = hostCommunity
			return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
		}

		if !checkPossibleOptions(config.Community, communityOptions) {
			logFields["parameter"] = hostCommunity
			return serror.New(fmt.Errorf(inCorrectValueOfParameter, config.Community, communityOptions), logFields)
		}
	} else {
		//check required fields for SNMP v3
		if !checkSetParameter(config.SecurityLevel) {
			logFields["parameter"] = hostSecurityLevel
			return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
		}

		if !checkPossibleOptions(config.SecurityLevel, securityLevelOptions) {
			logFields["parameter"] = hostSecurityLevel
			return serror.New(fmt.Errorf(inCorrectValueOfParameter, config.SecurityLevel, securityLevelOptions), logFields)
		}

		if !checkSetParameter(config.AuthProtocol) {
			logFields["parameter"] = hostAuthProtocol
			return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
		}

		if !checkPossibleOptions(config.AuthProtocol, authProtocolOptions) {
			logFields["parameter"] = hostAuthProtocol
			return serror.New(fmt.Errorf(inCorrectValueOfParameter, config.AuthProtocol, authProtocolOptions), logFields)
		}

		if !checkSetParameter(config.PrivProtocol) {
			logFields["parameter"] = hostPrivProtocol
			return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
		}

		if !checkPossibleOptions(config.PrivProtocol, privProtocolOptions) {
			logFields["parameter"] = hostPrivProtocol
			return serror.New(fmt.Errorf(inCorrectValueOfParameter, config.PrivProtocol, privProtocolOptions), logFields)
		}

		//set default values
		if !checkSetParameter(config.Retries) {
			config.Retries = defaultRetries
		}

		if !checkSetParameter(config.Timeout) {
			config.Timeout = defaultTimeout
		}

	}
	return nil
}

//readMetricConfigFile reads metric configuration from file and decodes it to structures
func readMetricConfigFile(setFilePath string) (Metrics, serror.SnapError) {
	var config Metrics
	logFields := map[string]interface{}{}
	logFields["setfile_path"] = setFilePath

	setFileContent, err := cfgReader.ReadFile(setFilePath)
	logFields["setfile_content"] = setFileContent
	if err != nil {
		return config, serror.New(err, logFields)
	}

	if len(setFileContent) == 0 {
		return config, serror.New(fmt.Errorf("Metrics configuration file is empty"), logFields)
	}

	err = json.Unmarshal(setFileContent, &config)
	if err != nil {
		return config, serror.New(fmt.Errorf("Settings file cannot be unmarshalled, err: %s", err), logFields)
	}
	return config, nil
}

//validateMetricConfig validates configuration of metrics
func validateMetricConfig(config Metrics) serror.SnapError {
	logFields := map[string]interface{}{}
	logFields["metric_config"] = config

	for metricName, metricDefinition := range config {
		logFields["metric_name"] = metricName

		//check OID -  requirered parameter
		if !checkSetParameter(metricDefinition.Oid) {
			logFields["parameter"] = metricOid
			return serror.New(fmt.Errorf(missingRequiredParameter), logFields)
		}

		//set default mode option if empty
		if !checkSetParameter(metricDefinition.Mode) {
			metricDefinition.Mode = ModeSingle
		}

		//check possible options for mode parameter
		if !checkPossibleOptions(metricDefinition.Mode, modeOptions) {
			logFields["parameter"] = metricMode
			return serror.New(fmt.Errorf(inCorrectValueOfParameter, metricDefinition.Mode, modeOptions), logFields)
		}

		//validate prefix fields
		if serr := validatePrefixSuffixFields(metricDefinition.Prefix); serr != nil {
			logFields["parameter"] = metricPrefix
			serr.SetFields(logFields)
			return serr
		}

		//validate suffix fields
		if serr := validatePrefixSuffixFields(metricDefinition.Suffix); serr != nil {
			logFields["parameter"] = metricSuffix
			serr.SetFields(logFields)
			return serr
		}

		//set default value for scale if scale is not configured
		if !checkSetParameter(metricDefinition.Scale) {
			metricDefinition.Scale = 1.0
		}
	}
	return nil
}

//validatePrefixSuffixFields validates configuration and fields of prefix or suffix which is added for metric name
func validatePrefixSuffixFields(fields map[string]string) serror.SnapError {
	if fields != nil {
		source, ok := fields[PrefixSuffixSource]
		if !ok {
			return serror.New(fmt.Errorf("Cannot find `source` parameter in configuration"))
		}

		if !checkPossibleOptions(source, sourcePrefixSuffixOptions) {
			return serror.New(fmt.Errorf("Incorrect value of parameter (%s), possible options: %v", source, sourcePrefixSuffixOptions))
		}

		if source == SourceSNMP {
			_, ok := fields[PrefixSuffixOID]
			if !ok {
				return serror.New(fmt.Errorf("Cannot find `OID` parameter in configuration"))
			}
		} else {
			_, ok := fields[PrefixSuffixString]
			if !ok {
				return serror.New(fmt.Errorf("Cannot find `string` parameter in configuration"))
			}
		}
	}
	return nil
}

//checkRequiredParam checks if required parameter is set
func checkSetParameter(param interface{}) bool {
	switch param.(type) {
	case string:
		if param.(string) == "" {
			return false
		}
	case int:
		if param.(int) == 0 {
			return false
		}
	case float64:
		if param.(float64) == 0 {
			return false
		}
	default:
		return false
	}
	return true
}

//checkPossibleOptions checks if value of parameter is one of possible values in configuration
func checkPossibleOptions(param interface{}, paramOptions []interface{}) bool {
	for _, opt := range paramOptions {
		if opt == param {
			return true
		}
	}
	return false
}
