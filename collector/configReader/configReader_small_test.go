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

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	CORRECT_METRIC_CONFIG_1 = iota
	CORRECT_METRIC_CONFIG_2
	CORRECT_METRIC_CONFIG_3
	CORRECT_METRIC_CONFIG_4
	WRONG_METRIC_CONFIG_1
	WRONG_METRIC_CONFIG_2
	WRONG_METRIC_CONFIG_3
	WRONG_METRIC_CONFIG_4
	WRONG_METRIC_CONFIG_5
	WRONG_METRIC_CONFIG_6
	SETFILE_NOT_FOUND
	EMPTY_SETFILE
	WRONG_SETFILE
	CORRECT_AGENT_CONFIG_1
	CORRECT_AGENT_CONFIG_2
	CORRECT_AGENT_CONFIG_3
	WRONG_AGENT_CONFIG_1
	WRONG_AGENT_CONFIG_2
	WRONG_AGENT_CONFIG_3
	WRONG_AGENT_CONFIG_4
	WRONG_AGENT_CONFIG_5
	WRONG_AGENT_CONFIG_6
	WRONG_AGENT_CONFIG_7
	WRONG_AGENT_CONFIG_8
	WRONG_AGENT_CONFIG_9
	WRONG_AGENT_CONFIG_10
	WRONG_AGENT_CONFIG_11
	WRONG_AGENT_CONFIG_12
	WRONG_AGENT_CONFIG_13
	WRONG_AGENT_CONFIG_14
	EMPTY_AGENT_CONFIG
)

type metricsEntry struct {
	out []byte
	err error
}

func newMetricsConfig(b []byte, e error) metricsEntry {
	return metricsEntry{b, e}
}

var metricsConfigsTestTable = map[int]metricsEntry{
	CORRECT_METRIC_CONFIG_1: newMetricsConfig(json.Marshal(getCorrectConfig1())),
	CORRECT_METRIC_CONFIG_2: newMetricsConfig(json.Marshal(getCorrectConfig2())),
	CORRECT_METRIC_CONFIG_3: newMetricsConfig(json.Marshal(getCorrectConfig3())),
	CORRECT_METRIC_CONFIG_4: newMetricsConfig(json.Marshal(getCorrectConfig4())),
	WRONG_METRIC_CONFIG_1:   newMetricsConfig(json.Marshal(getWrongConfig1())),
	WRONG_METRIC_CONFIG_2:   newMetricsConfig(json.Marshal(getWrongConfig2())),
	WRONG_METRIC_CONFIG_3:   newMetricsConfig(json.Marshal(getWrongConfig3())),
	WRONG_METRIC_CONFIG_4:   newMetricsConfig(json.Marshal(getWrongConfig4())),
	WRONG_METRIC_CONFIG_5:   newMetricsConfig(json.Marshal(getWrongConfig5())),
	WRONG_METRIC_CONFIG_6:   newMetricsConfig(json.Marshal(getWrongConfig6())),
	SETFILE_NOT_FOUND:       newMetricsConfig(nil, errors.New("Setfile not found")),
	EMPTY_SETFILE:           newMetricsConfig(nil, nil),
	WRONG_SETFILE:           newMetricsConfig(json.Marshal(map[string]int{"Foo": 1, "Bar": 2})),
}

var agentConfigsTestTable = map[int]map[string]interface{}{
	CORRECT_AGENT_CONFIG_1: getCorrectAgentConfig1(),
	CORRECT_AGENT_CONFIG_2: getCorrectAgentConfig2(),
	CORRECT_AGENT_CONFIG_3: getCorrectAgentConfig3(),
	WRONG_AGENT_CONFIG_1:   getWrongAgentConfig1(),
	WRONG_AGENT_CONFIG_2:   getWrongAgentConfig2(),
	WRONG_AGENT_CONFIG_3:   getWrongAgentConfig3(),
	WRONG_AGENT_CONFIG_4:   getWrongAgentConfig4(),
	WRONG_AGENT_CONFIG_5:   getWrongAgentConfig5(),
	WRONG_AGENT_CONFIG_6:   getWrongAgentConfig6(),
	WRONG_AGENT_CONFIG_7:   getWrongAgentConfig7(),
	WRONG_AGENT_CONFIG_8:   getWrongAgentConfig8(),
	WRONG_AGENT_CONFIG_9:   getWrongAgentConfig9(),
	WRONG_AGENT_CONFIG_10:  getWrongAgentConfig10(),
	WRONG_AGENT_CONFIG_11:  getWrongAgentConfig11(),
	WRONG_AGENT_CONFIG_12:  getWrongAgentConfig12(),
	WRONG_AGENT_CONFIG_13:  getWrongAgentConfig13(),
	WRONG_AGENT_CONFIG_14:  getWrongAgentConfig14(),
	EMPTY_AGENT_CONFIG:     getWrongAgentConfig16(),
}

type mockReader struct {
	vals metricsEntry
}

func (r *mockReader) ReadFile(s string) ([]byte, error) {
	return r.vals.out, r.vals.err
}

func TestGetMetricsConfig(t *testing.T) {
	Convey("Testing GetMetricsConfig", t, func() {

		Convey("Testing CORRECT_METRIC_CONFIG_1", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[CORRECT_METRIC_CONFIG_1]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldBeNil)
		})

		Convey("Testing CORRECT_METRIC_CONFIG_2", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[CORRECT_METRIC_CONFIG_2]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldBeNil)
		})

		Convey("Testing CORRECT_METRIC_CONFIG_3", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[CORRECT_METRIC_CONFIG_3]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldBeNil)
		})

		Convey("Testing CORRECT_METRIC_CONFIG_4", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[CORRECT_METRIC_CONFIG_4]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_1", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_1]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_2", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_2]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_3", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_3]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_4", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_4]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_5", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_5]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_METRIC_CONFIG_6", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_METRIC_CONFIG_6]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing SETFILE_NOT_FOUND", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[SETFILE_NOT_FOUND]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing EMPTY_SETFILE", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[EMPTY_SETFILE]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_SETFILE", func() {
			cfgReader = &mockReader{metricsConfigsTestTable[WRONG_SETFILE]}
			_, serr := GetMetricsConfig("setfile.json")
			So(serr, ShouldNotBeNil)
		})
	})
}

func TestGetAgentConfig(t *testing.T) {
	Convey("Testing GetAgentConfig", t, func() {

		Convey("Testing CORRECT_AGENT_CONFIG_1", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[CORRECT_AGENT_CONFIG_1])
			So(serr, ShouldBeNil)
		})

		Convey("Testing CORRECT_AGENT_CONFIG_2", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[CORRECT_AGENT_CONFIG_2])
			So(serr, ShouldBeNil)
		})

		Convey("Testing CORRECT_AGENT_CONFIG_3", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[CORRECT_AGENT_CONFIG_3])
			So(serr, ShouldBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_1", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_1])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_2", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_2])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_3", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_3])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_4", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_4])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_5", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_5])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_6", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_6])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_7", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_7])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_8", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_8])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_9", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_9])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_10", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_10])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_11", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_11])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_12", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_12])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_13", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_13])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing WRONG_AGENT_CONFIG_14", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[WRONG_AGENT_CONFIG_14])
			So(serr, ShouldNotBeNil)
		})

		Convey("Testing EMPTY_AGENT_CONFIG", func() {
			_, serr := GetSnmpAgentConfig(agentConfigsTestTable[EMPTY_AGENT_CONFIG])
			So(serr, ShouldNotBeNil)
		})

	})
}
