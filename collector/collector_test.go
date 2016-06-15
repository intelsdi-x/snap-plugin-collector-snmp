// +build legacy

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
	"os"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

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

		Convey("successfully obtain metrics name", func() {
			plg := New()
			createMockFile(mockFileCont)
			defer deleteMockFile()

			config := plugin.NewPluginConfigType()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})

			So(func() { plg.GetMetricTypes(config) }, ShouldNotPanic)
			_, err := plg.GetMetricTypes(config)
			So(err, ShouldBeNil)
		})
	})
}

func TestNew(t *testing.T) {
	Convey("Creating new plugin", t, func() {
		plugin := New()
		So(plugin, ShouldNotBeNil)
		So(plugin.metricConfig, ShouldNotBeNil)
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

		Convey("when no configuration settings available", func() {
			// set metrics config
			config := cdata.NewNode()
			mts := mockMts
			for i := range mts {
				mts[i].Config_ = config
			}

			plg := New()
			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			_, err := plg.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
		})

		Convey("when setfile configuration variable has incorrect type", func() {
			// set metrics config
			config := cdata.NewNode()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueInt{Value: 1})
			mts := mockMts
			for i := range mts {
				mts[i].Config_ = config
			}

			plg := New()
			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			_, err := plg.CollectMetrics(mts)
			So(err, ShouldNotBeNil)
		})

		Convey("when configuration is invalid", func() {
			//set metrics config
			config := cdata.NewNode()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})
			mts := mockMts
			for i := range mts {
				mts[i].Config_ = config
			}

			Convey("incorrect path to setfile", func() {
				// setfile does not exist
				deleteMockFile()

				plg := New()
				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				_, err := plg.CollectMetrics(mts)
				So(err, ShouldNotBeNil)
			})

			Convey("setfile is empty", func() {
				//setfile is empty
				createMockFile(mockFileContEmpty)
				defer deleteMockFile()

				plg := New()
				So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
				_, err := plg.CollectMetrics(mts)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("collect metrics successfully", func() {
			//create setfile
			createMockFile(mockFileCont)
			defer deleteMockFile()

			//set metrics config
			config := cdata.NewNode()
			config.AddItem(setFileConfigVar, ctypes.ConfigValueStr{Value: mockFilePath})
			config.AddItem("snmp_version", ctypes.ConfigValueStr{Value: "v2c"})
			config.AddItem("snmp_host_address", ctypes.ConfigValueStr{Value: "127.0.0.1"})
			config.AddItem("community", ctypes.ConfigValueStr{Value: "public"})
			mts := mockMts
			for i := range mts {
				mts[i].Config_ = config
			}

			plg := New()
			So(func() { plg.CollectMetrics(mts) }, ShouldNotPanic)
			_, err := plg.CollectMetrics(mts)
			So(err, ShouldBeNil)
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

	mockFileCont = []byte(`{
	    "__metric1__": {
	      "OID": ".1.3.6.1.2.1.1.9.1.3.1"
	    },
	  "__metric2__": {
	      "OID": ".1.3.6.1.2.1.1.9.1.3",
	      "mode": "multiple",
	      "prefix": {
	        "source": "snmp",
	        "OID": ".1.3.6.1.2.1.1.9.1.3"
	      },
	      "suffix": {
	        "source": "snmp",
	        "OID": ".1.3.6.1.2.1.1.9.1.3"
	      },
	      "scale": 1.25,
	      "shift": 2,
	      "unit": "",
	      "description": "description"
	    },
	  "__metric3__": {
	      "OID": ".1.3.6.1.2.1.1.9.1.3",
	      "mode": "multiple",
	      "prefix": {
	        "source": "string",
	        "string": "prefix"
	      },
	      "suffix": {
	        "source": "string",
	        "string": "suffix"
	      },
	      "scale": 1.25,
	      "shift": 2,
	      "unit": "",
	      "description": "description"
	    }
	}
	`)

	mockFileContEmpty = []byte(``)
)
