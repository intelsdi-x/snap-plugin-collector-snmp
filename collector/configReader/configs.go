// +build small

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

func getCorrectConfig1() Metrics {
	metricConfig := []Metric{Metric{
		Oid:  ".1.3.6.1.2.1.1.2.0",
		Mode: "single",
		Namespace: []Namespace{Namespace{Source: "string", String: "test1"},
			Namespace{Source: "string", String: "test2"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1", Description: "description"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name2", Description: "description"},
			Namespace{Source: "index", OidPart: 9, Name: "name3", Description: "description"},
			Namespace{Source: "string", String: "value"},
		},
		Unit:        "unit1",
		Description: "description1",
		Shift:       4.5,
		Scale:       1.2,
	}}
	return metricConfig
}

func getCorrectConfig2() Metrics {
	metricConfig := []Metric{Metric{
		Oid:  ".1.3.6.1.2.1.1.9.1.3",
		Mode: "walk",
		Namespace: []Namespace{
			Namespace{Source: "string", String: "test1"},
			Namespace{Source: "string", String: "test2"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1", Description: "description"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name2", Description: "description"},
			Namespace{Source: "index", OidPart: 9, Name: "name3", Description: "description"},
			Namespace{Source: "string", String: "value"},
		},
		Unit:        "unit2",
		Description: "description2",
		Shift:       -4.5,
		Scale:       -1.2,
	}}
	return metricConfig
}

func getCorrectConfig3() Metrics {
	metricConfig := []Metric{Metric{
		Oid:  ".1.3.6.1.2.1.1.2.0",
		Mode: "table",
		Namespace: []Namespace{
			Namespace{Source: "string", String: "test1"},
			Namespace{Source: "string", String: "test2"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1", Description: "description"},
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name2", Description: "description"},
			Namespace{Source: "index", OidPart: 9, Name: "name3", Description: "description"},
			Namespace{Source: "string", String: "value"},
		},
		Unit:        "unit3",
		Description: "description3",
	}}
	return metricConfig
}

func getCorrectConfig4() Metrics {
	metricConfig := []Metric{Metric{
		Oid: ".1.3.6.1.2.1.1.9.1.3",
		Namespace: []Namespace{
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1", Description: "description"},
			Namespace{Source: "string", String: "test2"}},
	}}
	return metricConfig
}

func getWrongConfig1() Metrics {
	metricConfig := []Metric{Metric{}}
	return metricConfig
}

func getWrongConfig2() Metrics {
	metricConfig := []Metric{Metric{
		Oid:       ".1.3.6.1.2.1.1.2.0",
		Namespace: []Namespace{Namespace{Source: "string"}},
	}}
	return metricConfig
}

func getWrongConfig3() Metrics {
	metricConfig := []Metric{Metric{
		Oid: ".1.3.6.1.2.1.1.2.0",
		Namespace: []Namespace{
			Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Description: "description"},
		},
		Unit:        "unit",
		Description: "description",
	}}
	return metricConfig
}

func getWrongConfig4() Metrics {
	metricConfig := []Metric{Metric{
		Oid: ".1.3.6.1.2.1.1.2.0",
		Namespace: []Namespace{
			Namespace{Source: "index", OidPart: 9, Description: "description"},
		},
	}}
	return metricConfig
}

func getWrongConfig5() Metrics {
	metricConfig := []Metric{Metric{
		Oid: ".1.3.6.1.2.1.1.2.0",
		Namespace: []Namespace{
			Namespace{Source: "snmp", Name: "name1", Description: "description"},
		},
		Unit:        "unit",
		Description: "description",
	}}
	return metricConfig
}

func getWrongConfig6() Metrics {
	metricConfig := []Metric{Metric{
		Oid: ".1.3.6.1.2.1.1.2.0",
		Namespace: []Namespace{
			Namespace{Source: "index", Name: "name3", Description: "description"},
		},
		Unit:        "unit",
		Description: "description",
	}}
	return metricConfig
}

func getCorrectAgentConfig1() map[string]interface{} {
	//configuration for SNMP v1 and SNMP v2c
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v2c"
	agentConfig["community"] = "public"
	return agentConfig
}

func getCorrectAgentConfig2() map[string]interface{} {
	//configuration for SNMP v3
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}

func getCorrectAgentConfig3() map[string]interface{} {
	//all possible options
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	agentConfig["security_engine_id"] = ""
	agentConfig["context_engine_id"] = ""
	agentConfig["context_engine_id"] = ""
	agentConfig["context_name"] = ""
	agentConfig["retries"] = 3
	agentConfig["timeout"] = 5
	return agentConfig
}

func getWrongAgentConfig1() map[string]interface{} {
	//missing required parameter snmp_version
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["community"] = "public"
	return agentConfig
}

func getWrongAgentConfig2() map[string]interface{} {
	//incorrect value of snmp_version parameter
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "incorrectVersion"
	agentConfig["community"] = "public"
	return agentConfig
}
func getWrongAgentConfig3() map[string]interface{} {
	//incorrect type of value of snmp_version parameter
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = 2
	agentConfig["community"] = "public"
	return agentConfig
}
func getWrongAgentConfig4() map[string]interface{} {
	//incorrect type of value of community
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v2c"
	agentConfig["community"] = 0
	return nil
}
func getWrongAgentConfig5() map[string]interface{} {
	//incorrect security level
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "incorrectSecurityLevel"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}
func getWrongAgentConfig6() map[string]interface{} {
	//incorrect auth_protocol
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "IncorrectAuthProtocol"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}
func getWrongAgentConfig7() map[string]interface{} {
	//incorrect priv_protocol
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "InocorrectPrivProtocol"
	return agentConfig
}

func getWrongAgentConfig8() map[string]interface{} {
	//unsupported type of snmp version
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = true
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}
func getWrongAgentConfig9() map[string]interface{} {
	//missing required parameter -  community
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v2c"
	return agentConfig
}

func getWrongAgentConfig10() map[string]interface{} {
	//missing required parameter - security_level
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}

func getWrongAgentConfig11() map[string]interface{} {
	//missing required parameter - auth_protocol
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}

func getWrongAgentConfig12() map[string]interface{} {
	//missing required parameter - auth_protocol
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}

func getWrongAgentConfig13() map[string]interface{} {
	//missing required parameter - priv_protocol
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_agent_address"] = "127.0.0.1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	return agentConfig
}

func getWrongAgentConfig14() map[string]interface{} {
	//missing required parameter - snmp_agent_address
	agentConfig := make(map[string]interface{})
	agentConfig["snmp_agent_name"] = "agent1"
	agentConfig["snmp_version"] = "v3"
	agentConfig["user_name"] = "user"
	agentConfig["security_level"] = "NoAuthNoPriv"
	agentConfig["auth_password"] = "password"
	agentConfig["auth_protocol"] = "MD5"
	agentConfig["priv_password"] = "password"
	agentConfig["priv_protocol"] = "DES"
	return agentConfig
}

func getWrongAgentConfig16() map[string]interface{} {
	//empty agent configuration
	return nil
}
