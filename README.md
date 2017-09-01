# snap-plugin-collector-snmp

This plugin collects metrics using SNMP (Simple Network Management Protocol). The plugin sends GET and GETNEXT requests to receive metrics from SNMP agents using numeric OID (Object Identifier) which indicates metric value or set of values.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Snap's Global Config](#snaps-global-config)
  * [Setfile structure](#setfile-structure)
  * [Namespace](#namespace)
  * [Metric modes](#modes)
  * [SNMP agent configuration](#snmp-agent-configuration)
  * [Task Manifest](#task-manifest)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements
* [golang 1.6+](https://golang.org/dl/) - needed only for building
* access to SNMP agent which supports SNMP in one of the following versions: v1, v2c, v3

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation

#### Download snmp plugin binary:

You can get the pre-built binaries for your OS and architecture at plugin's [Github Releases](https://github.com/intelsdi-x/snap-plugin-collector-snmp/releases) page.

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-snmp

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-snmp
```
Build the plugin by running `make` in the repo:
```
$ make
```
This builds the plugin in `./build`.

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).

* Create configuration file (called a *Setfile*) in which metrics are defined. *Setfile* structure description is available in [the setfile structure section](#setfile-structure) and there are examples in [examples/setfiles/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/setfiles/).

* Create Global Config, see description in [Snap's Global Config](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/README.md#snaps-global-config).

* Create a [Task Manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) with SNMP agent configuration (described in [this](#snmp-agent-configuration)). There are examples in the [Task Manifest](#task-manifest) section or in [examples/tasks/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/).
 
Notice that this plugin is a generic plugin, it cannot work without configuration, because there is no reasonable default behavior.

## Documentation

### Collected Metrics
The plugin collects metrics using SNMP.

Metrics are available in namespaces which are configurable by the user. Namespaces start with `/intel/snmp/`, further parts of namespaces need to be configured, details are described in sections: [setfile structure](#setfile-structure) and [namespace](#namespace).
For each of metrics following tags are added:
- OID - object identifier which is used to read metric,
- SNMP_AGENT_NAME - name given by the user for SNMP agent in configuration of SNMP agent,
- SNMP_AGENT_ADDRESS - IP address or host name with port number of SNMP agent.

Metric names are defined in *Setfile* and can be collected in one of following data types: int32, uint32, uint64, float64, string. 

Detailed descriptions of data types are available in the table below:

SNMP data type | SNMP plugin data type | Description
----------------|:-------------------------|:-----------------------
 Counter, Counter32| uint32 | Represents a non-negative integer which monotonically increases until it reaches a maximum value of 32bits-1 (4294967295 dec), when it wraps around and starts increasing again from zero
 Counter64 | uint64 |Same as Counter32 but has a maximum value of 64bits-1
 Gauge32  | uint32 | Represents an unsigned integer, which may increase or decrease, but shall never exceed a maximum value
 Integer | int32 | Signed 32bit Integer (values between -2147483648 and 2147483647)
 Integer32 | int32 | Same as Integer
 IpAddress | string | IP address
 Object Identifier | string | An OID
 Octet String | string | Arbitrary binary or textual data, typically limited to 255 characters in length
 TimeTicks | uint32 | Represents an unsigned integer which represents the time, modulo 232 (4294967296 dec), in hundredths of a second between two epochs
 UInteger32  | uint32 | Unsigned 32bit Integer (values between 0 and 4294967295)
 
 
#### Modification of metric value

It is possible to modify metric value using `scale` or `shift` parameters, for more information read [setfile structure](#setfile-structure) section.

The metric value is modified using the following equation:

`new_metric_value = numeric_metric_value * scale + shift`.

If `scale` or `shift` parameters are set (`scale` different than 1, `shift` different than 0) then numeric metrics are returned as float64.

### snap's Global Config
Global configuration files are described in [Snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTELD_CONFIGURATION.md) and require the `snmp` section in `collector`
along with the specific *Setfile* - path to SNMP plugin configuration file (path to *Setfile*).
Examples of valid Global Config files are in [examples/cfg/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/configs/).

It is useful to set higher value of `max_running_plugins` in global configuration because, for SNMP plugin, for each of tasks a one instance of plugin is needed.
Default value of `max_running_plugins` is 3 so by default only 3 tasks with SNMP plugin can be created.

### Setfile structure

Setfile contains JSON structure which is used to define metrics. Each metric is defined as JSON object in the following format:
```
    {
        "namespace": {
            {"source": "string", "string": "<string>"},
            {"source": "snmp", "OID": "<object_identifier>", "name": "<name>", "description": "<description>"},
            {"source": "index", "oid_part": <oid_part_number>, "name": "<name>", "description": "<description>"},
        }
      "OID": "<object_identifier>",
      "mode": "<metric_mode>",
      "scale": <scale_value>,
      "shift": <shift_value>,
      "unit": "<unit>",
      "description": "<description>"
    }
```
Detailed descriptions of all parameters in metric definition are available in the table below:

Parameter | Type | Possible options | Required | Description
----------------|:-------------------------|:-----------------------|:-----------------------|:-----------------------
 namespace |  array  | - | yes | Array of configuration for namespace elements
 namespace::source | string | string/snmp/index |  yes | Source of namespace element, namespace elements can be defined as string value (*string*), can be received using SNMP request (*snmp*), or can be defined as a number from OID (*index*), see [namespace section](#namespace)
 namespace::string | string | - | yes, for source set to *string* | Namespace element defined by the user as a string value
 namespace::OID | string | - | yes, for source set to *snmp* | Numeric OID which is used to receive namespace element
 namespace::oid_part | uint | - | yes, for source set to *index* | Index of OID part which is used in namespace. It indicates part of OID which will be used in namespace, counting parts (numbers in OID) of OID from 0
 namespace::name | string | - | yes, for source set to *index* or *snmp* | Name of dynamic metric
 namespace::description | string | - | yes, for source set to *index* or *snmp* | Description of dynamic metric
 OID  | string | - | yes | Object identifier
 mode | string | single/table/walk | no | Mode of metric, it is possible to read a single metric or read metrics from the specific node of MIB (ang. Management Information Base), see [metric modes section](#modes), on default *single* is set
 unit |  string | - | no | Metric unit
 description | string | - | no | Metric description
 shift | float64 | - | no | Shift value can be added to numeric metric
 scale | float64 | - | no | Numeric metric can be multiplied by scale value


Here is an example metric definition (with more available in [examples/setfiles/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/setfiles/)):
```
[
  {
  "mode": "single",
  "namespace": [
    {"source": "string", "string": "net-single"},
    {"source": "string", "string": "if-single"},
    {"source": "snmp", "name": "interface", "description": "interface name", "OID": ".1.3.6.1.2.1.2.2.1.2.1"},
    {"source": "index", "name": "id", "description": "number from OID", "oid_part": 10},
    {"source": "string", "string": "in_octets"}
  ],
  "OID": ".1.3.6.1.2.1.2.2.1.10.1",
  "scale": 1.0,
  "shift": 0,
  "unit": "unit",
  "description": "description of metric"
  },
 {
  "mode": "table",
  "namespace": [
    {"source": "string", "string": "net-table"},
    {"source": "string", "string": "if-table"},
    {"source": "snmp", "name": "interface",  "description": "interface name", "OID": ".1.3.6.1.2.1.2.2.1.2"},
    {"source": "index", "name": "id", "description": "number from OID", "oid_part": 10},
    {"source": "string", "string": "in_octets"}
  ],
  "OID": ".1.3.6.1.2.1.2.2.1.10",
  "scale": 1.0,
  "shift": 0,
  "unit": "unit",
  "description": "description  of metric"
 },
 {
  "mode": "walk",
  "namespace": [
    {"source": "string", "string": "net-walk"},
    {"source": "string", "string": "if-walk"},
    {"source": "index", "name": "index9", "description": "number from OID", "oid_part": 9},
    {"source": "index", "name": "index10", "description": "number from OID", "oid_part": 10},
    {"source": "string", "string": "value"}
  ],
  "OID": ".1.3.6.1.2.1.2.2.1.",
  "scale": 1.0,
  "shift": 0,
  "unit": "unit",
  "description": "description of metric"
}
]
```

### Namespace

Metrics namespaces are configured in *Setfile*. Namespaces start with `/intel/snmp/`,  further parts of namespaces need to be configured. 

Namespace is configured as an array which contains configuration of namespace elements (separted by `/`). For example metrics definition are shown in [setfile structure section](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/README.md#setfile-structure) and this MIB:

```
.1.3.6.1.2.1.2.1.0 = INTEGER: 2
.1.3.6.1.2.1.2.2.1.1.1 = INTEGER: 1
.1.3.6.1.2.1.2.2.1.1.2 = INTEGER: 2
.1.3.6.1.2.1.2.2.1.2.1 = STRING: lo
.1.3.6.1.2.1.2.2.1.2.2 = STRING: eth0
.1.3.6.1.2.1.2.2.1.3.1 = INTEGER: softwareLoopback(24)
.1.3.6.1.2.1.2.2.1.3.2 = INTEGER: ethernetCsmacd(6)
.1.3.6.1.2.1.2.2.1.4.1 = INTEGER: 65536
.1.3.6.1.2.1.2.2.1.4.2 = INTEGER: 9001
.1.3.6.1.2.1.2.2.1.5.1 = Gauge32: 10000000
.1.3.6.1.2.1.2.2.1.5.2 = Gauge32: 4294967295
.1.3.6.1.2.1.2.2.1.6.1 = STRING: 
.1.3.6.1.2.1.2.2.1.6.2 = STRING: 
.1.3.6.1.2.1.2.2.1.7.1 = INTEGER: up(1)
.1.3.6.1.2.1.2.2.1.7.2 = INTEGER: up(1)
.1.3.6.1.2.1.2.2.1.8.1 = INTEGER: up(1)
.1.3.6.1.2.1.2.2.1.8.2 = INTEGER: up(1)
.1.3.6.1.2.1.2.2.1.9.1 = Timeticks: (0) 0:00:00.00
.1.3.6.1.2.1.2.2.1.9.2 = Timeticks: (0) 0:00:00.00
.1.3.6.1.2.1.2.2.1.10.1 = Counter32: 2426340104
.1.3.6.1.2.1.2.2.1.10.2 = Counter32: 2116292856
```

following namespaces are built:
- for *single* mode:
```
/intel/snmp/net-single/if-single/lo/1/in_octets
```
- for *table* mode:
```
/intel/snmp/net-table/if-table/lo/1/in_octets
/intel/snmp/net-table/if-table/eth0/2/in_octets
```
- for *walk* mode:
```
/intel/snmp/net-walk/if-walk/1/1/value
/intel/snmp/net-walk/if-walk/1/2/value
/intel/snmp/net-walk/if-walk/2/1/value
/intel/snmp/net-walk/if-walk/2/2/value
/intel/snmp/net-walk/if-walk/3/1/value
/intel/snmp/net-walk/if-walk/3/2/value
/intel/snmp/net-walk/if-walk/4/1/value
/intel/snmp/net-walk/if-walk/4/2/value
/intel/snmp/net-walk/if-walk/5/1/value
/intel/snmp/net-walk/if-walk/5/2/value
/intel/snmp/net-walk/if-walk/6/1/value
/intel/snmp/net-walk/if-walk/6/2/value
/intel/snmp/net-walk/if-walk/7/1/value
/intel/snmp/net-walk/if-walk/7/2/value
/intel/snmp/net-walk/if-walk/8/1/value
/intel/snmp/net-walk/if-walk/8/2/value
/intel/snmp/net-walk/if-walk/9/1/value
/intel/snmp/net-walk/if-walk/9/2/value
/intel/snmp/net-walk/if-walk/10/1/value
/intel/snmp/net-walk/if-walk/10/2/value
```
Length of namespace can be different but the last element in array must have *source* option set to *string*.

### Metric modes

There are three modes to gather SNMP metrics:

- `single` - mode to read only one metric
- `table` - mode to read set of metrics from one node
- `walk` - mode to read set of metrics from multiple nodes, all children nodes are read

### SNMP agent configuration

SNMP agent configuration is created in Task Manifest, in the `config` section `/intel/snmp` section must be created and set of appropriate SNMP agent parameters must be configured. All possible parameters for SNMP agent are gathered in the table below:

Parameter | Type | Possible options | Valid for SNMP  versions | Default value | Required | Description
----------------|:-------------------------|:-----------------------|:-----------------------|:-----------------------|:-----------------------|:-----------------------
 snmp_agent_name | string | - |v1,v2c,v3 | -  | no | SNMP agent name give by the user, any string helpful for the user, this parameter is added as tag (SNMP_AGENT_NAME) for metrics
 snmp_agent_address | string | - | v1,v2c,v3 | - | yes | IP address or host name with port number. This parameter is added as a tag (SNMP_AGENT_ADDRESS) for metrics
 snmp_version | string | v1/v2c/v3 | v1,v2c,v3 | -  | yes | SNMP version
 community | string | - | v1,v2c | - | yes | Community
 user_name | string | - |  v3 |  - | yes | User name
 security_level | string | NoAuthNoPriv/AuthNoPriv/AuthPriv | v3 | - | yes | Security level
 auth_password | string | - | v3 | - | no |  Authentication protocol pass phrase
 auth_protocol | string | MD5/SHA | v3 | - | yes | Authentication protocol
 priv_password | string | - | v3 | - | no | Privacy protocol pass phrase
 priv_protocol | string | DES/AES| v3 | - | yes | Privacy protocol
 security_engine_id | string| - | v3 | - | no | Security engine ID
 context_engine_id | string | - | v3 | - | no | Context engine ID
 context_name | string | - | v3 | - | no | Context name 
 retries | uint | - | v1,v2c,v3 | 1 | no | Number of connection retries 
 timeout | int | -  | v1,v2c,v3 | 5 | no | SNMP request timeout in seconds
 
 *WARNING:* Notice that `retries` and `timeout` and also `interval` in Task Manifest files must be adjusted to SNMP agent responsiveness. Unsuitable values of these parameters could cause problems with metrics collection (some metrics could be missing).
 
### Task Manifest

Example [Task Manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (more examples in [examples/tasks/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/)):
```
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "30s"
  },
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/snmp/*": {}
      },
      "config": {
        "/intel/snmp": {
          "snmp_agent_name": "host1",
          "snmp_agent_address": "127.0.0.1:1161",
          "snmp_version": "v2c",
          "community": "public",
          "network": "tcp",
          "timeout": 5,
          "retries": 5
        }
      },
      "publish": [
        {
          "plugin_name": "file",
          "config": {
            "file": "/tmp/published_snmp.txt"
          }
        }
      ]
    }
  }
}
```

### Examples

Here's an example of running the snap-plugin-collector-snmp plugin and writing data to a file using [snap-plugin-publisher-file](https://github.com/intelsdi-x/snap-plugin-publisher-file).

- Create configuration file (*Setfile*) for SNMP plugin based on the examples in [examples/setfiles/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/setfiles/)

- Set path to configuration file as a field in the *Setfile* in Global Config based on examples in [examples/configs/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/configs/)

- Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started),
in one terminal window, run `snapteld` (in this case with logging set to 1, trust disabled and global configuration saved in config.json ):
```
$ snapteld -l 1 -t 0 --config config.json
```

- In another terminal window, download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-snmp/latest/linux/x86_64/snap-plugin-collector-snmp
$ snaptel plugin load snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-snmp
```

See available metrics for your system:

```
$ snaptel metric list
```

- Create a [Task Manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) file to use snap-plugin-collector-snmp plugin based on the files in [examples/tasks/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/))

- Create a task:
```
$ snaptel task create -t task.json
```

- And watch the metrics populate:
```
$ snaptel task watch <task_id>
```

- To stop previously created task:
```
$ snaptel task stop <task_id>
```

#### Running in plugin diagnostic mode:
Plugin can be started in plugin diagnostic mode without need for Snap daemon to be running. This mode can be used for checking plugin output in development process.

General rule for runnig diagnostic mode is to launch plugin binary. This plugin needs configuration, so it is needed to pass that configuration as argument. For example:

```bash
$ ./build/linux/x86_64/snap-plugin-collector-snmp --config '{
    "setfile": "setfile.json",
    "snmp_agent_name": "host1",
    "snmp_agent_address": "127.0.0.1:161",
    "snmp_version": "v3",
    "network": "udp",
    "user_name": "user",
    "security_level": "AuthPriv",
    "auth_password": "password",
    "auth_protocol": "MD5",
    "priv_protocol": "DES",
    "priv_password": "password"
}'
```

If configuration is valid, plugin should output metric catalog and collected metrics to standard output.
As runnig diagnostic mode command for this plugin is not handy, you can find above example as Bash script in `examples/standalone.sh`.

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-users/issues) and feel free to then submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-users/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, the open telemetry framework. See the full project at http://github.com/intelsdi-x/snap. To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Katarzyna Kujawa](https://github.com/katarzyna-z)

