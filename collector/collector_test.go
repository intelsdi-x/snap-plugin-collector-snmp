// +build medium

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
	"os"
	"testing"

	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/configReader"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/intelsdi-x/snap/core/serror"
	"github.com/k-sone/snmpgo"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	SUCCESSFULLY_CREATED_HANDLER = iota
	UNSUCCESSFULLY_CREATED_HANDLER
	SNMP_ELEMENT_CORRECT_OCTET_STRING
	SNMP_ELEMENT_CORRECT_COUNTER32
	SNMP_ELEMENT_CORRECT_COUNTER64
	SNMP_ELEMENT_CORRECT_INTEGER
	SNMP_ELEMENT_INCORRECT
)

type snmpHandlerEntry struct {
	s   *snmpgo.SNMP
	err serror.SnapError
}

func newSNMPHandler(s *snmpgo.SNMP, e serror.SnapError) snmpHandlerEntry {
	return snmpHandlerEntry{s, e}
}

type snmpMock struct {
	handlerEntry snmpHandlerEntry
	elementEntry snmpElementEntry
}

var snmpHandlerTestTable = map[int]snmpHandlerEntry{
	SUCCESSFULLY_CREATED_HANDLER:   newSNMPHandler(&snmpgo.SNMP{}, nil),
	UNSUCCESSFULLY_CREATED_HANDLER: newSNMPHandler(&snmpgo.SNMP{}, serror.New(fmt.Errorf("Error - new handler not created"))),
}

func (m *snmpMock) newHandler(hostConfig configReader.SnmpAgent) (*snmpgo.SNMP, serror.SnapError) {
	return m.handlerEntry.s, m.handlerEntry.err
}

func (m *snmpMock) readElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, serror.SnapError) {
	varBinds := []*snmpgo.VarBind{m.elementEntry.element}
	return varBinds, m.elementEntry.err
}

type snmpElementEntry struct {
	element *snmpgo.VarBind
	err     serror.SnapError
}

func newElementEntry(oid string, val snmpgo.Variable, e serror.SnapError) snmpElementEntry {
	newOid, _ := snmpgo.NewOid(oid)
	return snmpElementEntry{snmpgo.NewVarBind(newOid, val), e}
}

var snmpElementTestTable = map[int]snmpElementEntry{
	SNMP_ELEMENT_CORRECT_OCTET_STRING: newElementEntry(".1.3.6.1.2.1.1.9.1.3.1", snmpgo.NewOctetString([]byte("variable123")), nil),
	SNMP_ELEMENT_CORRECT_COUNTER32:    newElementEntry(".1.3.6.1.2.1.1.9.1.3.1", snmpgo.NewCounter32(123), nil),
	SNMP_ELEMENT_CORRECT_COUNTER64:    newElementEntry(".1.3.6.1.2.1.1.9.1.3.1", snmpgo.NewCounter64(123), nil),
	SNMP_ELEMENT_CORRECT_INTEGER:      newElementEntry(".1.3.6.1.2.1.1.9.1.3.1", snmpgo.NewInteger(123), nil),
	SNMP_ELEMENT_INCORRECT:            newElementEntry("", nil, serror.New(fmt.Errorf("Error - snmp request fails"))),
}

func (m *snmpMock) ReadElement(handler *snmpgo.SNMP, oid string) (*snmpgo.VarBind, serror.SnapError) {
	return m.elementEntry.element, m.elementEntry.err
}

func TestMeta(t *testing.T) {
	Convey("Calling Meta function", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, pluginType)
		So(meta.AcceptedContentTypes[0], ShouldEqual, plugin.SnapGOBContentType)
		So(meta.ReturnedContentTypes[0], ShouldEqual, plugin.SnapGOBContentType)
		So(meta.RoutingStrategy, ShouldEqual, plugin.StickyRouting)
	})
}

func TestGetMetricTypes(t *testing.T) {
	Convey("Getting exposed metric types", t, func() {

		Convey("when no configuration item available", func() {
			cfg := plugin.NewPluginConfigType()
			plugin := New()
			So(func() { plugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			_, err := plugin.GetMetricTypes(cfg)
			So(err, ShouldNotBeNil)
		})

		Convey("when path to setfile is incorrect", func() {
			plg := New()
			// file has not existed yet
			deleteMockFile()

			//create configuration
			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			So(func() { plg.GetMetricTypes(config) }, ShouldNotPanic)
			_, err := plg.GetMetricTypes(config)
			So(err, ShouldNotBeNil)

		})

		Convey("when setfle configuration variable has incorrect type (int instead of string)", func() {
			plg := New()
			// file has not existed yet
			deleteMockFile()

			//create configuration
			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueInt{Value: 3})

			So(func() { plg.GetMetricTypes(config) }, ShouldNotPanic)
			_, err := plg.GetMetricTypes(config)
			So(err, ShouldNotBeNil)

		})

		Convey("when setfile is empty", func() {

			plg := New()
			createMockFile(mockFileContEmpty)
			defer deleteMockFile()

			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			So(func() { plg.GetMetricTypes(config) }, ShouldNotPanic)
			_, err := plg.GetMetricTypes(config)
			So(err, ShouldNotBeNil)
		})

		Convey("when setfile content is incorrect - wrong `oid_part` parametr", func() {
			plg := New()
			createMockFile(mockFileContWrong)
			defer deleteMockFile()

			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			_, err := plg.GetMetricTypes(config)

			So(err, ShouldNotBeNil)
		})

		Convey("successfully obtain metrics name", func() {
			plg := New()
			createMockFile(mockFileCont)
			defer deleteMockFile()

			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			mts, err := plg.GetMetricTypes(config)

			So(err, ShouldBeNil)
			So(len(mts), ShouldEqual, 7)

			Convey("then correct list of metrics is returned", func() {
				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, m.Namespace().String())
				}
				So(namespaces, ShouldContain, "/intel/snmp/hrSystem/*/*/*/value")
				So(namespaces, ShouldContain, "/intel/snmp/system/sysORTable/sysOREntry/sysORDescr/*/value")
				So(namespaces, ShouldContain, "/intel/snmp/hostName")
				So(namespaces, ShouldContain, "/intel/snmp/*/*/name")
				So(namespaces, ShouldContain, "/intel/snmp/hrSystemNumUsers")
				So(namespaces, ShouldContain, "/intel/snmp/hrSystemProcesses")
				So(namespaces, ShouldContain, "/intel/snmp/sysServicesModified")
			})
		})
	})
}

func TestNew(t *testing.T) {
	Convey("Creating new plugin", t, func() {
		plugin := New()
		So(plugin, ShouldNotBeNil)
		So(plugin.metricsConfigs, ShouldNotBeNil)
	})
}

func TestGetConfigPolicy(t *testing.T) {
	plugin := New()

	Convey("Getting config policy", t, func() {
		So(func() { plugin.GetConfigPolicy() }, ShouldNotPanic)
		configPolicy, err := plugin.GetConfigPolicy()
		So(err, ShouldBeNil)
		So(configPolicy, ShouldNotBeNil)
	})
}

func TestCollectMetrics(t *testing.T) {
	Convey("Collecting metrics", t, func() {

		Convey("when collect metrics successfully", func() {

			plg := New()

			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			//create plugin config
			pluginConfig := plugin.NewPluginConfigType()
			pluginConfig.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			//setfile is read in GetMetricsTypes
			mts, err := plg.GetMetricTypes(pluginConfig)
			So(err, ShouldBeNil)

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}

			Convey("when received data with OCTET_STRING type", func() {
				snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
					elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				metrics, err := plg.CollectMetrics(mts)

				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeEmpty)
				So(metrics[0].Data().(string), ShouldEqual, "variable123")
			})

			Convey("when received data with COUNTER32  type", func() {
				snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
					elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_COUNTER32]}

				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				metrics, err := plg.CollectMetrics(mts)

				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeEmpty)
				So(metrics[0].Data().(uint64), ShouldEqual, 123)
			})

			Convey("when received data with COUNTER64  type", func() {
				snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
					elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_COUNTER64]}

				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				metrics, err := plg.CollectMetrics(mts)

				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeEmpty)
				So(metrics[0].Data().(uint64), ShouldEqual, 123)
			})

			Convey("when received data with INTEGER  type", func() {
				snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
					elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_INTEGER]}

				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				metrics, err := plg.CollectMetrics(mts)

				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeEmpty)
				So(metrics[0].Data().(int64), ShouldEqual, 123)
			})

			Convey("when metric have dynamic namespace elements", func() {
				snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
					elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_INTEGER]}

				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				metrics, err := plg.CollectMetrics(mts)

				So(err, ShouldBeNil)
				So(metrics, ShouldNotBeEmpty)
				isDynamic, dynamics := metrics[0].Namespace().IsDynamic()
				So(isDynamic, ShouldEqual, true)
				So(dynamics, ShouldHaveLength, 3)
				So(dynamics, ShouldContain, 3)
				So(dynamics, ShouldContain, 4)
				So(dynamics, ShouldContain, 5)
			})
		})

		Convey("when setfile content is incorrect - wrong `oid_part` parametr", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//create setfile
			createMockFile(mockFileContWrong)
			defer deleteMockFile()

			plg := New()

			//create metric
			mts := []plugin.MetricType{plugin.MetricType{Namespace_: core.NewNamespace(vendor, pluginName, "__metric1__")}}

			//create plugin config
			pluginConfig := plugin.NewPluginConfigType()
			pluginConfig.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}
			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)

			//force initialization
			plg.initialized = false

			_, err := plg.CollectMetrics(mts)

			So(err, ShouldNotBeNil)
		})

		Convey("when metric name is incorrect", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			plg := New()

			//create metric
			mts := []plugin.MetricType{plugin.MetricType{Namespace_: core.NewNamespace(vendor, pluginName, "incorrect_metric")}}

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}

			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			mts, err := plg.CollectMetrics(mts)

			So(err, ShouldNotBeNil)
			So(mts, ShouldBeEmpty)
		})

		Convey("when cannot create a new snmp handler", func() {
			//clear connections map
			snmpConnections = make(map[string]connection)

			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[UNSUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			plg := New()

			//create metric
			mts := []plugin.MetricType{plugin.MetricType{Namespace_: core.NewNamespace(vendor, pluginName, "__metric1__")}}

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}

			mts, err := plg.CollectMetrics(mts)

			So(err, ShouldNotBeNil)
			So(mts, ShouldBeEmpty)
		})

		Convey("when host configuration is incorrect", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[UNSUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			plg := New()

			//create metric w namespace with correct namespace length
			mts := []plugin.MetricType{plugin.MetricType{Namespace_: core.NewNamespace(vendor, pluginName, "__metric1__")}}

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "incorrect snmp version"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}

			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			mts, err := plg.CollectMetrics(mts)

			So(err, ShouldNotBeNil)
			So(mts, ShouldBeEmpty)
		})

		Convey("when snmp request fails", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_INCORRECT]}

			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			plg := New()

			//create plugin config
			pluginConfig := plugin.NewPluginConfigType()
			pluginConfig.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			//setfile is read in GetMetricsTypes
			mts, err := plg.GetMetricTypes(pluginConfig)
			So(err, ShouldBeNil)

			//create host config
			config := cdata.NewNode()
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_agent_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			for i := range mts {
				mts[i].Config_ = config
			}

			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			mts, err = plg.CollectMetrics(mts)

			So(err, ShouldBeNil)
			So(mts, ShouldBeEmpty)
		})

	})
}

func TestConvertSnmpDataToMetric(t *testing.T) {

	Convey("Calling convertSnmpDataToMetric ", t, func() {

		Convey("with parsing error for Counter", func() {
			_, serr := convertSnmpDataToMetric("12.3", "Counter")
			So(serr, ShouldNotBeNil)
		})

		Convey("with parsing error for Counter64", func() {
			_, serr := convertSnmpDataToMetric("12.3", "Counter64")
			So(serr, ShouldNotBeNil)
		})

		Convey("with parsing error for Integer", func() {
			_, serr := convertSnmpDataToMetric("12.3", "Integer")
			So(serr, ShouldNotBeNil)
		})

		Convey("with parsing git derror for unsupported type", func() {
			_, serr := convertSnmpDataToMetric("12.3", "Integer64")
			So(serr, ShouldBeNil)
		})

	})
}

func TestGetDynamicNamespaceElements(t *testing.T) {
	Convey("Calling getDynamicNamespaceElements ", t, func() {

		//create SNMP agent config
		agentConfig := make(map[string]interface{})
		agentConfig["snmp_agent_name"] = "agent1"
		agentConfig["snmp_agent_address"] = "127.0.0.1"
		agentConfig["snmp_version"] = "v2c"
		agentConfig["community"] = "public"
		snmpAgentConfig, _ := configReader.GetSnmpAgentConfig(agentConfig)

		Convey("with incorrect value of `oid_part` parameter", func() {

			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//metric configuration
			metricConfig := []configReader.Metric{configReader.Metric{
				Oid:  ".1.3.6.1.2.1.1.2.0",
				Mode: "table",
				Namespace: []configReader.Namespace{
					configReader.Namespace{Source: "string", String: "test1"},
					configReader.Namespace{Source: "string", String: "test2"},
					configReader.Namespace{Source: "index", OidPart: 102, Name: "name"},
					configReader.Namespace{Source: "string", String: "value"},
				},
				Unit:        "unit",
				Description: "description",
			}}

			//create results
			newOid, err := snmpgo.NewOid(metricConfig[0].Oid)
			So(err, ShouldBeNil)

			varBind := snmpgo.NewVarBind(newOid, snmpgo.NewCounter32(123))
			varBinds := []*snmpgo.VarBind{varBind}

			handler, err := snmp_.newHandler(snmpAgentConfig)
			So(err, ShouldBeNil)

			serr := getDynamicNamespaceElements(handler, varBinds, &metricConfig[0])
			So(serr, ShouldNotBeNil)

		})

		Convey("with incorrect configuration of dynamic elements of namespace, the number of results is not equal the number of dynamic elements of namespace ", func() {

			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//metric configuration
			metricConfig := []configReader.Metric{configReader.Metric{
				Oid:  ".1.3.6.1.2.1.1.2.0",
				Mode: "table",
				Namespace: []configReader.Namespace{
					configReader.Namespace{Source: "string", String: "test1"},
					configReader.Namespace{Source: "string", String: "test2"},
					configReader.Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name"},
					configReader.Namespace{Source: "string", String: "value"},
				},
				Unit:        "unit",
				Description: "description",
			}}

			//create results
			newOid, err := snmpgo.NewOid(metricConfig[0].Oid)
			So(err, ShouldBeNil)

			varBind := snmpgo.NewVarBind(newOid, snmpgo.NewCounter32(123))
			varBinds := []*snmpgo.VarBind{varBind, varBind}

			handler, err := snmp_.newHandler(snmpAgentConfig)
			So(err, ShouldBeNil)

			serr := getDynamicNamespaceElements(handler, varBinds, &metricConfig[0])
			So(serr, ShouldNotBeNil)

		})

		Convey("when SNMP request fails", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_INCORRECT]}

			//metric configuration
			metricConfig := []configReader.Metric{configReader.Metric{
				Oid:  ".1.3.6.1.2.1.1.2.0",
				Mode: "table",
				Namespace: []configReader.Namespace{
					configReader.Namespace{Source: "string", String: "test1"},
					configReader.Namespace{Source: "string", String: "test2"},
					configReader.Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1"},
					configReader.Namespace{Source: "string", String: "value"},
				},
				Unit:        "unit",
				Description: "description",
			}}

			//create results
			newOid, err := snmpgo.NewOid(metricConfig[0].Oid)
			So(err, ShouldBeNil)

			varBind := snmpgo.NewVarBind(newOid, snmpgo.NewCounter32(123))
			varBinds := []*snmpgo.VarBind{varBind}

			handler, err := snmp_.newHandler(snmpAgentConfig)
			So(err, ShouldBeNil)

			serr := getDynamicNamespaceElements(handler, varBinds, &metricConfig[0])
			So(serr, ShouldNotBeNil)
		})

		Convey("when dynamic elements of namespace are successfully received", func() {
			snmp_ = &snmpMock{handlerEntry: snmpHandlerTestTable[SUCCESSFULLY_CREATED_HANDLER],
				elementEntry: snmpElementTestTable[SNMP_ELEMENT_CORRECT_OCTET_STRING]}

			//metric configuration
			metricConfig := []configReader.Metric{configReader.Metric{
				Oid:  ".1.3.6.1.2.1.1.2.0",
				Mode: "table",
				Namespace: []configReader.Namespace{
					configReader.Namespace{Source: "string", String: "test1"},
					configReader.Namespace{Source: "string", String: "test2"},
					configReader.Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1"},
					configReader.Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name2"},
					configReader.Namespace{Source: "index", OidPart: 3, Name: "name3"},
					configReader.Namespace{Source: "string", String: "value"},
				},
				Unit:        "unit",
				Description: "description",
			}}

			//create results
			newOid, err := snmpgo.NewOid(metricConfig[0].Oid)
			So(err, ShouldBeNil)

			varBind := snmpgo.NewVarBind(newOid, snmpgo.NewCounter32(123))
			varBinds := []*snmpgo.VarBind{varBind}

			handler, err := snmp_.newHandler(snmpAgentConfig)
			So(err, ShouldBeNil)

			serr := getDynamicNamespaceElements(handler, varBinds, &metricConfig[0])
			So(serr, ShouldBeNil)
		})
	})
}

func TestGetMetricsToCollect(t *testing.T) {
	Convey("Calling getMetricsToCollect ", t, func() {
		metricConfig := configReader.Metric{
			Oid:  ".1.3.6.1.2.1.1.2.0",
			Mode: "table",
			Namespace: []configReader.Namespace{
				configReader.Namespace{Source: "string", String: "test1"},
				configReader.Namespace{Source: "string", String: "test2"},
				configReader.Namespace{Source: "snmp", Oid: ".1.3.6.1.2.1.1.9.1.3", Name: "name1"},
				configReader.Namespace{Source: "string", String: "value"},
			},
			Unit:        "unit",
			Description: "description",
		}

		metricsConfigs := map[string]configReader.Metric{
			"/intel/snmp/test1/test2/*/value": metricConfig,
			"/intel/snmp/test1/test3/*/value": metricConfig,
		}

		Convey("with correct arguments", func() {

			collectedMetrics, serr := getMetricsToCollect("/intel/snmp/test1/test2/*/value", metricsConfigs)

			So(serr, ShouldBeNil)
			So(len(collectedMetrics), ShouldEqual, 1)
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Oid, ShouldEqual, ".1.3.6.1.2.1.1.2.0")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Mode, ShouldEqual, "table")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Unit, ShouldEqual, "unit")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Description, ShouldEqual, "description")
			So(len(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace), ShouldEqual, 4)
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[0].String, ShouldEqual, "test1")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[0].Source, ShouldEqual, "string")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[1].String, ShouldEqual, "test2")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[1].Source, ShouldEqual, "string")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[2].Source, ShouldEqual, "snmp")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[2].Oid, ShouldEqual, ".1.3.6.1.2.1.1.9.1.3")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[3].String, ShouldEqual, "value")
			So(collectedMetrics["/intel/snmp/test1/test2/*/value"].Namespace[3].Source, ShouldEqual, "string")
		})

		Convey("with incorrect regex expression", func() {

			collectedMetrics, serr := getMetricsToCollect("a(b", metricsConfigs)

			So(serr, ShouldNotBeNil)
			So(len(collectedMetrics), ShouldEqual, 0)

		})
	})
}

func createMockFile(fileCont []byte) {
	deleteMockFile()

	f, _ := os.Create(mockFilePath)
	f.Write(fileCont)
}

func deleteMockFile() {
	os.Remove(mockFilePath)
}

var (
	mockMts = []plugin.MetricType{
		plugin.MetricType{Namespace_: core.NewNamespace(vendor, pluginName, "*", "value")},
	}

	mockFilePath = "./temp_setfile.json"

	mockFileCont = []byte(`
		[
		 {
		  "mode": "walk",
		  "namespace": [
			{"source": "string", "string": "hrSystem"},
			{"source": "index", "name": "index9", "description": "number from OID", "oid_part": 8},
			{"source": "index", "name": "index10", "description": "number from OID", "oid_part": 9},
			{"source": "snmp", "name": "snmp part", "description": "part received through SNMP request", "OID": ".1.3.6.1.2.1.25.1"},
			{"source": "string", "string": "value"}
		  ],
		  "OID": ".1.3.6.1.2.1.25.1",
		  "scale": 1.0,
		  "shift": 0,
		  "unit": "unit",
		  "description": "description of metric"
		},
		   {
		  "mode": "table",
		  "namespace": [
			{"source": "string", "string": "system"},
			{"source": "string", "string": "sysORTable"},
			{"source": "string", "string": "sysOREntry"},
			{"source": "string", "string": "sysORDescr"},
			{"source": "snmp", "name": "snmp part", "description": "part received through SNMP request", "OID": ".1.3.6.1.2.1.1.9.1.3"},
			{"source": "string", "string": "value"}
		  ],
		  "OID": ".1.3.6.1.2.1.1.9.1.4",
		  "scale": 1.0,
		  "shift": 0,
		  "unit": "unit",
		  "description": "description of metric"
		},
		   {
			"mode": "single",
			"namespace": [
			  {"source": "string", "string": "hostName"}
			],
			"OID": ".1.3.6.1.2.1.1.5.0",
			"description": "host name"
		  },
		   {
			"mode": "single",
			"namespace": [
			  {"source": "snmp", "name": "e-mail", "description": "part for email", "OID": ".1.3.6.1.2.1.1.4.0"},
			  {"source": "index", "name": "index name", "description": "part for index", "oid_part": 8},
			  {"source": "string", "string": "name"}
			],
			"OID": ".1.3.6.1.2.1.1.5.0",
			"description": "host name - with dynamic namespace"
		  },
		  {
			"mode": "single",
			"namespace": [
			  {"source": "string", "string": "hrSystemNumUsers"}
			],
			"OID": ".1.3.6.1.2.1.25.1.5.0",
			"scale": 1,
			"shift": 0,
			"unit": "",
			"description": "Numeric metric"
		  },
		   {
			"mode": "single",
			"namespace": [
			  {"source": "string", "string": "hrSystemProcesses"}
			],
			"OID": ".1.3.6.1.2.1.25.1.6.0",
			"scale": 1,
			"shift": 0,
			"unit": "",
			"description": "Numeric metric"
		  },
		  {
			"mode": "single",
			"namespace": [
			  {"source": "string", "string": "sysServicesModified"}
			],
			"OID": ".1.3.6.1.2.1.1.7.0",
			"scale": 0.5,
			"shift": 18.5,
			"unit": "",
			"description": "Numeric metric to show usage of scale and shift"
		  },
		   {
			"mode": "single",
			"namespace": [
			  {"source": "string", "string": "sysServicesModified"}
			],
			"OID": ".1.3.6.1.2.1.1.7.0",
			"scale": 0.5,
			"shift": 18.5,
			"unit": "",
			"description": "Numeric metric to show usage of scale and shift"
		  }
	]
 `)

	mockFileContWrong = []byte(`
		[
		 {
		  "mode": "walk",
		  "namespace": [
			{"source": "string", "string": "hrSystem"},
			{"source": "index", "name": "index9", "description": "number from OID", "oid_part": -8},
			{"source": "index", "name": "index10", "description": "number from OID", "oid_part": 9},
			{"source": "snmp", "name": "snmp part", "description": "part received through SNMP request", "OID": ".1.3.6.1.2.1.25.1"},
			{"source": "string", "string": "value"}
		  ],
		  "OID": ".1.3.6.1.2.1.25.1",
		  "scale": 1.0,
		  "shift": 0,
		  "unit": "unit",
		  "description": "description of metric"
		},
		 {
		  "mode": "walk",
		  "namespace": [
			{"source": "string", "string": "hrSystem"},
			{"source": "index", "name": "index9", "description": "number from OID", "oid_part": "abc"},
			{"source": "index", "name": "index10", "description": "number from OID", "oid_part": 9},
			{"source": "snmp", "name": "snmp part", "description": "part received through SNMP request", "OID": ".1.3.6.1.2.1.25.1"},
			{"source": "string", "string": "value"}
		  ],
		  "OID": ".1.3.6.1.2.1.25.1",
		  "scale": 1.0,
		  "shift": 0,
		  "unit": "unit",
		  "description": "description of metric"
		}
	]
 `)

	mockFileContEmpty = []byte(``)
)
