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
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:         ".1.3.6.1.2.1.1.2.0",
		Mode:        "single",
		Prefix:      map[string]string{"source": "string", "string": "prefix_"},
		Suffix:      map[string]string{"source": "string", "string": "_suffix"},
		Unit:        "unit1",
		Description: "description1",
		Shift:       4.5,
		Scale:       1.2,
	}
	return metricConfig
}

func getCorrectConfig2() Metrics {
	pluginConfig1 := make(Metrics)
	pluginConfig1["metric2"] = Metric{
		Oid:         ".1.3.6.1.2.1.1.9.1.3",
		Mode:        "multiple",
		Prefix:      map[string]string{"source": "string", "string": "prefix_"},
		Suffix:      map[string]string{"source": "string", "string": "_suffix"},
		Unit:        "unit2",
		Description: "description2",
		Shift:       -4.5,
		Scale:       -1.2,
	}
	return pluginConfig1
}

func getCorrectConfig3() Metrics {
	pluginConfig1 := make(Metrics)
	pluginConfig1["metric3"] = Metric{
		Oid:         ".1.3.6.1.2.1.1.2.0",
		Mode:        "single",
		Prefix:      map[string]string{"source": "snmp", "OID": ".1.3.6.1.2.1.1.2.0"},
		Suffix:      map[string]string{"source": "snmp", "OID": ".1.3.6.1.2.1.1.2.0"},
		Unit:        "unit3",
		Description: "description3",
	}
	return pluginConfig1
}

func getCorrectConfig4() Metrics {
	pluginConfig1 := make(Metrics)
	pluginConfig1["metric4"] = Metric{
		Oid:         ".1.3.6.1.2.1.1.9.1.3",
		Mode:        "multiple",
		Prefix:      map[string]string{"source": "snmp", "OID": ".1.3.6.1.2.1.1.9.1.3"},
		Suffix:      map[string]string{"source": "snmp", "OID": ".1.3.6.1.2.1.1.9.1.3"},
		Unit:        "unit4",
		Description: "description4",
	}
	return pluginConfig1
}

func getWrongConfig1() Metrics {
	pluginConfig1 := make(Metrics)
	pluginConfig1["metric1"] = Metric{}
	return pluginConfig1
}

func getWrongConfig2() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Prefix: map[string]string{"source": "string", "xxx": ""},
	}
	return metricConfig
}

func getWrongConfig3() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Prefix: map[string]string{"source": "snmp", "xxx": ""},
	}
	return metricConfig
}

func getWrongConfig4() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:  ".1.3.6.1.2.1.1.2.0",
		Mode: "incorrectMode",
	}
	return metricConfig
}

func getWrongConfig5() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Prefix: map[string]string{"xxx": ""},
	}
	return metricConfig
}

func getWrongConfig6() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Suffix: map[string]string{"source": "string", "xxx": ""},
	}
	return metricConfig
}

func getWrongConfig7() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Suffix: map[string]string{"source": "snmp", "xxx": ""},
	}
	return metricConfig
}

func getWrongConfig8() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Prefix: map[string]string{"source": "incorrectPrefixSource"},
	}
	return metricConfig
}

func getWrongConfig9() Metrics {
	metricConfig := make(Metrics)
	metricConfig["metric1"] = Metric{
		Oid:    ".1.3.6.1.2.1.1.2.0",
		Suffix: map[string]string{"source": "incorrectSuffixSource"},
	}
	return metricConfig
}

func getCorrectHostConfig1() map[string]interface{} {
	//configuration for SNMP v1 and SNMP v2c
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v2c"
	hostConfig["community"] = "public"
	return hostConfig
}

func getCorrectHostConfig2() map[string]interface{} {
	//configuration for SNMP v3
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}

func getCorrectHostConfig3() map[string]interface{} {
	//all possible options
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	hostConfig["security_engine_id"] = ""
	hostConfig["context_engine_id"] = ""
	hostConfig["context_engine_id"] = ""
	hostConfig["context_name"] = ""
	hostConfig["retries"] = 3
	hostConfig["timeout"] = 5
	return hostConfig
}

func getWrongHostConfig1() map[string]interface{} {
	//missing required parameter snmp_version
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["community"] = "public"
	return hostConfig
}

func getWrongHostConfig2() map[string]interface{} {
	//incorrect value of snmp_version parameter
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "incorrectVersion"
	hostConfig["community"] = "public"
	return hostConfig
}
func getWrongHostConfig3() map[string]interface{} {
	//incorrect type of value of snmp_version parameter
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = 2
	hostConfig["community"] = "public"
	return hostConfig
}
func getWrongHostConfig4() map[string]interface{} {
	//incorrect type of value of community
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v2c"
	hostConfig["community"] = 0
	return nil
}
func getWrongHostConfig5() map[string]interface{} {
	//incorrect security level
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "incorrectSecurityLevel"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}
func getWrongHostConfig6() map[string]interface{} {
	//incorrect auth_protocol
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "IncorrectAuthProtocol"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}
func getWrongHostConfig7() map[string]interface{} {
	//incorrect priv_protocol
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "InocorrectPrivProtocol"
	return hostConfig
}

func getWrongHostConfig8() map[string]interface{} {
	//unsupported type of snmp version
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = true
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}
func getWrongHostConfig9() map[string]interface{} {
	//missing required parameter -  community
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v2c"
	return hostConfig
}

func getWrongHostConfig10() map[string]interface{} {
	//missing required parameter - security_level
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}

func getWrongHostConfig11() map[string]interface{} {
	//missing required parameter - auth_protocol
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}

func getWrongHostConfig12() map[string]interface{} {
	//missing required parameter - auth_protocol
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}

func getWrongHostConfig13() map[string]interface{} {
	//missing required parameter - priv_protocol
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_host_address"] = "127.0.0.1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	return hostConfig
}

func getWrongHostConfig14() map[string]interface{} {
	//missing required parameter - snmp_host_address
	hostConfig := make(map[string]interface{})
	hostConfig["snmp_host_name"] = "host1"
	hostConfig["snmp_version"] = "v3"
	hostConfig["user_name"] = "user"
	hostConfig["security_level"] = "NoAuthNoPriv"
	hostConfig["auth_password"] = "password"
	hostConfig["auth_protocol"] = "MD5"
	hostConfig["priv_password"] = "password"
	hostConfig["priv_protocol"] = "DES"
	return hostConfig
}

func getWrongHostConfig16() map[string]interface{} {
	//empty host configuration
	return nil
}
