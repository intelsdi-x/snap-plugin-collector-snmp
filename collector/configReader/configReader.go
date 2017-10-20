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

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

const (
	//ModeSingle option in mode of metric
	ModeSingle = "single"

	//ModeWalk option in mode of metric
	ModeWalk = "walk"

	//ModeTable option in mode of metric
	ModeTable = "table"

	//nsSourceSNMP option in source of namespace element configuration
	NsSourceSNMP = "snmp"

	//nsSourceString option in source of namespace element configuration
	NsSourceString = "string"

	//nsSourceIndex option in source of namespace element configuration
	NsSourceIndex = "index"

	//agentName indicates SNMP agent name
	agentName = "snmp_agent_name"

	//agentAddress indicates SNMP agent address
	agentAddress = "snmp_agent_address"

	//agentSnmpVersion indicates SNMP version in SNMP agent configuration
	agentSnmpVersion = "snmp_version"

	//agentCommunity indicates community (SNMP  v1 &  SNMP v2c) in SNMP agent configuration
	agentCommunity = "community"

	//agentNetwork indicates network which is used in SNMP agent configuration, see net.Dial parameter
	agentNetwork = "network"

	//agentUserName indicates user name (SNMP v3) in SNMP agent configuration
	agentUserName = "user_name"

	//agentSecurityLevel indicates security level (SNMP v3) in SNMP agent configuration
	agentSecurityLevel = "security_level"

	//agentAuthPassword indicates authentication protocol pass phrase (SNMP v3) in SNMP agent configuration
	agentAuthPassword = "auth_password"

	//agentAuthProtocol indicates authentication protocol (SNMP v3) in SNMP agent configuration
	agentAuthProtocol = "auth_protocol"

	//agentPrivPassword indicates privacy protocol pass phrase (SNMP v3) in SNMP agent configuration
	agentPrivPassword = "priv_password"

	//agentPrivProtocol indicates privacy protocol (SNMP v3) in SNMP agent configuration
	agentPrivProtocol = "priv_protocol"

	//agentSecurityEngineId indicates security engine ID (SNMP v3) in SNMP agent configuration
	agentSecurityEngineId = "security_engine_id"

	//agentContextEngineID indicates context engine ID (SNMP v3) in SNMP agent configuration
	agentContextEngineID = "context_engine_id"

	//agentContextName indicates context name (SNMP v3) in SNMP agent configuration
	agentContextName = "context_name"

	//agentRetries indicates number of connection retries in SNMP agent configuration
	agentRetries = "retries"

	//agentTimeout indicates timeout for network connection in SNMP agent configuration
	agentTimeout = "timeout"

	//metricNamespace indicates metric namespace
	metricNamespace = "namespace"

	//metricOid indicates OID which is use to receive metric
	metricOid = "OID"

	//metricMode indicates mode of metric
	metricMode = "mode"

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
	missingRequiredParameter = "Missing required parameter in configuration (%s)"

	//incorrectValueOfParameter error message for incorrect value of parameter
	incorrectValueOfParameter = "Incorrect value of parameter (%s), possible options: %v"
)

type SnmpAgent struct {
	Name             string `mapstructure:"snmp_agent_name"`
	SnmpVersion      string `mapstructure:"snmp_version"`
	Address          string `mapstructure:"snmp_agent_address"`
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

type Namespace struct {
	Source      string `json:"source"`
	Name        string `json:"name"`
	String      string `json:"string"`
	OidPart     uint   `json:"oid_part"`
	Oid         string `json:"OID"`
	Description string `json:"description"`
	Values      []string
}

type Metric struct {
	Mode        string      `json:"mode"`
	Namespace   []Namespace `json:"namespace"`
	Oid         string      `json:"OID"`
	Unit        string      `json:"unit"`
	Description string      `json:"description"`
	Shift       float64     `json:"shift"`
	Scale       float64     `json:"scale"`
}

type Metrics []Metric

type cfgReaderType struct{}

type reader interface {
	ReadFile(s string) ([]byte, error)
}

var (
	//AgentConfigParameters slice of agent configuration parameters
	SnmpAgentConfigParameters = []string{agentName, agentAddress, agentSnmpVersion, agentCommunity, agentNetwork,
		agentUserName, agentSecurityLevel, agentAuthPassword, agentAuthProtocol, agentPrivPassword,
		agentPrivProtocol, agentSecurityEngineId, agentContextEngineID, agentContextName, agentRetries, agentTimeout}

	//modeOptions slice of options for mode parameter
	modeOptions = []interface{}{ModeSingle, ModeWalk, ModeTable}

	//snmpVersionOptions slice of options for SNMP version
	snmpVersionOptions = []interface{}{snmpv1, snmpv2, snmpv3}

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

//GetMetricsConfig decodes and validates configuration of SNMP agent
func GetSnmpAgentConfig(configMap plugin.Config) (SnmpAgent, error) {
	config, err := decodeSnmpAgentConfig(configMap)
	if err != nil {
		return config, err
	}

	err = validateSnmpAgentConfig(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

//GetMetricsConfig reads and validates configuration of metrics
func GetMetricsConfig(setFilePath string) (Metrics, error) {
	config, err := readMetricConfigFile(setFilePath)
	if err != nil {
		return config, err
	}

	err = validateMetricConfig(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

//decodeSnmpAgentConfig decodes configuration of SNMP agent into structure
func decodeSnmpAgentConfig(config plugin.Config) (SnmpAgent, error) {
	var snmpAgentConfig SnmpAgent
	logFields := map[string]interface{}{}
	err := mapstructure.Decode(config, &snmpAgentConfig)
	if err != nil {
		log.WithFields(logFields).Warn(err)
		return snmpAgentConfig, err
	}
	return snmpAgentConfig, nil
}

//validateSnmpAgentConfig validates configuration of SNMP agent
func validateSnmpAgentConfig(config SnmpAgent) error {
	logFields := map[string]interface{}{}
	logFields["agent_config"] = config

	if !checkSetParameter(config.Address) {
		logFields["parameter"] = agentAddress
		err := fmt.Errorf(missingRequiredParameter, agentAddress)
		log.WithFields(logFields).Warn(err)
		return err
	}

	if !checkSetParameter(config.SnmpVersion) {
		logFields["parameter"] = agentSnmpVersion
		err := fmt.Errorf(missingRequiredParameter, agentSnmpVersion)
		log.WithFields(logFields).Warn(err)
		return err
	}

	if !checkPossibleOptions(config.SnmpVersion, snmpVersionOptions) {
		logFields["parameter"] = agentSnmpVersion
		err := fmt.Errorf(incorrectValueOfParameter, config.SnmpVersion, snmpVersionOptions)
		log.WithFields(logFields).Warn(err)
		return err
	}

	if config.SnmpVersion == snmpv1 || config.SnmpVersion == snmpv2 {
		//check required fields for SNMP v1 and SNMP v2c
		if !checkSetParameter(config.Community) {
			logFields["paramgiteter"] = agentCommunity
			err := fmt.Errorf(missingRequiredParameter, agentCommunity)
			log.WithFields(logFields).Warn(err)
			return err
		}
	} else {
		//check required fields for SNMP v3
		if !checkSetParameter(config.SecurityLevel) {
			logFields["parameter"] = agentSecurityLevel
			err := fmt.Errorf(missingRequiredParameter, agentSecurityLevel)
			log.WithFields(logFields).Warn(err)
			return err
		}

		if !checkPossibleOptions(config.SecurityLevel, securityLevelOptions) {
			logFields["parameter"] = agentSecurityLevel
			err := fmt.Errorf(incorrectValueOfParameter, config.SecurityLevel, securityLevelOptions)
			log.WithFields(logFields).Warn(err)
			return err
		}

		if !checkSetParameter(config.AuthProtocol) {
			logFields["parameter"] = agentAuthProtocol
			err := fmt.Errorf(missingRequiredParameter, agentAuthProtocol)
			log.WithFields(logFields).Warn(err)
			return err
		}

		if !checkPossibleOptions(config.AuthProtocol, authProtocolOptions) {
			logFields["parameter"] = agentAuthProtocol
			err := fmt.Errorf(incorrectValueOfParameter, config.AuthProtocol, authProtocolOptions)
			log.WithFields(logFields).Warn(err)
			return err
		}

		if !checkSetParameter(config.PrivProtocol) {
			logFields["parameter"] = agentPrivProtocol
			err := fmt.Errorf(missingRequiredParameter, agentPrivProtocol)
			log.WithFields(logFields).Warn(err)
			return err
		}

		if !checkPossibleOptions(config.PrivProtocol, privProtocolOptions) {
			logFields["parameter"] = agentPrivProtocol
			err := fmt.Errorf(incorrectValueOfParameter, config.PrivProtocol, privProtocolOptions)
			log.WithFields(logFields).Warn(err)
			return err
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
func readMetricConfigFile(setFilePath string) (Metrics, error) {
	var config Metrics
	logFields := map[string]interface{}{}
	logFields["setfile_path"] = setFilePath

	setFileContent, err := cfgReader.ReadFile(setFilePath)
	logFields["setfile_content"] = setFileContent
	if err != nil {
		log.WithFields(logFields).Warn(err)
		return config, err
	}

	if len(setFileContent) == 0 {
		err := fmt.Errorf("Metrics configuration file is empty")
		log.WithFields(logFields).Warn(err)
		return config, err
	}

	err = json.Unmarshal(setFileContent, &config)
	if err != nil {
		err := fmt.Errorf("Settings file cannot be unmarshalled, err: %s", err)
		log.WithFields(logFields).Warn(err)
		return config, err
	}
	return config, nil
}

//validateMetricConfig validates configuration of metrics
func validateMetricConfig(metricConfigs Metrics) error {
	logFields := map[string]interface{}{}
	logFields["metric_config"] = metricConfigs

	for i := 0; i < len(metricConfigs); i++ {
		logFields["namespace_config"] = metricConfigs[i].Namespace

		//check namespace -  required parameter
		if !checkSetParameter(metricConfigs[i].Namespace) {
			logFields["parameter"] = metricNamespace
			err := fmt.Errorf(missingRequiredParameter, metricNamespace)
			log.WithFields(logFields).Warn(err)
			return err
		}
		//validate namespace configuration
		if err := validateNamespace(metricConfigs[i].Namespace); err != nil {
			logFields["parameter"] = metricNamespace
			log.WithFields(logFields).Warn(err)
			return err
		}

		lastNamespacePart := metricConfigs[i].Namespace[len(metricConfigs[i].Namespace)-1]
		if lastNamespacePart.Source != NsSourceString {
			logFields["parameter"] = metricNamespace
			err := fmt.Errorf("The last namespace element must have `source` set to `string`")
			log.WithFields(logFields).Warn(err)
			return err
		}

		//check OID -  required parameter
		if !checkSetParameter(metricConfigs[i].Oid) {
			logFields["parameter"] = metricOid
			err := fmt.Errorf(missingRequiredParameter, metricOid)
			log.WithFields(logFields).Warn(err)
			return err
		}

		//set default mode option if empty
		if !checkSetParameter(metricConfigs[i].Mode) {
			metricConfigs[i].Mode = ModeSingle
		}

		//check possible options for mode parameter
		if !checkPossibleOptions(metricConfigs[i].Mode, modeOptions) {
			logFields["parameter"] = metricMode
			err := fmt.Errorf(incorrectValueOfParameter, metricConfigs[i].Mode, modeOptions)
			log.WithFields(logFields).Warn(err)
			return err
		}

		//set default value for scale if scale is not configured
		if !checkSetParameter(metricConfigs[i].Scale) {
			metricConfigs[i].Scale = 1.0
		}
	}
	return nil
}

//validateNamespace validates configuration of metric namespace
func validateNamespace(namespaceConfig []Namespace) error {
	if len(namespaceConfig) == 0 {
		return fmt.Errorf("Missing namespace configuration")
	}

	for _, nsCfg := range namespaceConfig {
		switch nsCfg.Source {
		case NsSourceString:
			//check required  parameter for source set to string
			if !checkSetParameter(nsCfg.String) {
				return fmt.Errorf("Cannot find `string` parameter in configuration namespace element")
			}
		case NsSourceSNMP:
			//check required  parameter for source set to snmp
			if !checkSetParameter(nsCfg.Oid) {
				return fmt.Errorf("Cannot find `OID` parameter in configuration namespace element")
			}

			if !checkSetParameter(nsCfg.Name) {
				return fmt.Errorf("Cannot find `name` parameter in configuration namespace element")
			}

			if !checkSetParameter(nsCfg.Description) {
				return fmt.Errorf("Cannot find `description` parameter in configuration namespace element")
			}
		case NsSourceIndex:
			//check required  parameter for source set to index
			if !checkSetParameter(nsCfg.OidPart) {
				return fmt.Errorf("Cannot find `oid_part` parameter in configuration namespace element")
			}

			if !checkSetParameter(nsCfg.Name) {
				return fmt.Errorf("Cannot find `name` parameter in configuration namespace element")
			}

			if !checkSetParameter(nsCfg.Description) {
				return fmt.Errorf("Cannot find `description` parameter in configuration namespace element")
			}
		default:
			return fmt.Errorf("Incorrect value of `source` (%s) in namespace configuration, possible options: %v",
				nsCfg.Source, []string{NsSourceString, NsSourceSNMP, NsSourceIndex})
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
	case uint:
		if param.(uint) == 0 {
			return false
		}
	case float64:
		if param.(float64) == 0 {
			return false
		}
	case []Namespace:
		if param.([]Namespace) == nil {
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
