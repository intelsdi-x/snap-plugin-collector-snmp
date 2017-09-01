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

package main

import (
	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func main() {

	plg := collector.New()
	if plg == nil {
		panic("Plugin could not be initialized")
	}

	plugin.StartCollector(plg, collector.PluginName, collector.Version, plugin.RoutingStrategy(plugin.StickyRouter), plugin.ConcurrencyCount(1))
}
