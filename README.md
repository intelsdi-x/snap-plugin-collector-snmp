# snap-plugin-collector-snmp

This plugin collects metrics using SNMP (ang. Simple Network Management Protocol). The plugin sends GET and GETNEXT requests to receive metrics from SNMP agents.

The plugin is a generic plugin. You need to configure metrics.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [snap's Global Config](#snaps-global-config)
  * [Setfile structure](#setfile-structure)
  * [Task manifest](#task-manifest)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

* [golang 1.5+](https://golang.org/dl/) - needed only for building

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation

#### Download snmp plugin binary:

You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-snmp

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-snmp
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`.

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).

* Create configuration file (called as a *setfile*) in which metrics are defined, *setfile* structure description is available in [setfile structure section](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/README.md#setfile-structure) and see exemplary in [examples/setfiles/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/setfiles/).

* Create Global Config, see description in [snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/README.md#snaps-global-config).

* Create task manifest with SNMP host configuration, SNMP host configuration is described in [task manifest section](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/README.md#task-manifest) and see exemplary in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/).
 
Notice that this plugin is a generic plugin, it cannot work without configuration, because there is no reasonable default behavior.

## Documentation

### Collected Metrics
The plugin collects metrics using SNMP.

Metrics are available in namespace: `/intel/snmp/<metric_name>/value`, for each of metrics following tags are added:
- OID - object identifier which is used to read metric,                                                                 
- SNMP_HOST_NAME - name given by the user for SNMP host in configuration of host,
- SNMP_HOST_ADDRESS - IP address or host name with port number of SNMP host.

Metrics names are defined in setfile.
Metrics can be collected in one of following data types: int32, uint32, uint64, string. Detailed descriptions of data types are available in the table below:

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

### snap's Global Config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "snmp" in "collector" section and then specify `"setfile"` - path to snmp plugin configuration file (path to setfile),
see example Global Config in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/configs/).

### Setfile structure

Setfile contains JSON structure which is used to define metrics. 
Metric is defined as JSON object in following format:
```
	"<metric_name>": {
      "OID": "<object_identifier>",
      "mode": "<metric_mode>",
      "prefix": {
        "source": "<prefix_source>",
        "string": "<prefix>"
      },
      "suffix": {
        "source": "<suffix_source>",
        "string": "<OID_to_get_suffix>"
      },
      "scale": <scale_value>,
      "shift": <shift_value>,
      "unit": "<unit>",
      "description": "<description>"
    }
```

Detailed descriptions of all parameters in metric definition are available in the table below:

Parameter | Type | Possible options | Required | Description
----------------|:-------------------------|:-----------------------|:-----------------------|:-----------------------
 \<metric\_name\>  | string |  - | yes | Key in metrics' map.
 OID  | string | - | yes | Object identifier
 mode | string | single/multiple | no | Mode of metric, it is possible to read a single metric or read metrics (all elements) from the specific node of MIB (ang. Management Information Base)
 unit |  string | - | no | Metric unit
 description | string | - | no | Metric description
 prefix | object | - | no | Prefix for metric name, used in namespace
 prefix::source | string | string/snmp| no | Source of prefix, prefixcan be received either using SNMP request or it can be defined by user as string value
 prefix::string | string | - | no | Prefix for metric name provided by user
 prefix::OID | string | - | no | Numeric OID which should be used to receive prefix for metric name
 suffix | object | - | no | Suffix for metric name, used in namespace
 suffix::source | string | string/snmp| - | Source of suffix, suffix can be received either using SNMP request or it can be defined by user as string value
 suffix::string | string | - | no | Suffix for metric name provided by user
 suffix::OID | string | - | no | Numeric OID which should be used to receive suffix for metric name
 shift | float64 | - | no | Shift value can be added to numeric metric
 scale | float64 | - | no | Numeric metric can be multiplied by scale value


Exemplary metric definition:
```
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
      "unit": "unit",
      "description": "description"
    }
```
### Task manifest

SNMP host configuration is defined in task manifest, in the config section "/intel/snmp" section must be created and set of appropriate SNMP host parameters must be configured. All possible parameters for SNMP host are gathered in the table below:

Parameter | Type | Possible options | Valid for SNMP  versions | Default value | Required | Description
----------------|:-------------------------|:-----------------------|:-----------------------|:-----------------------|:-----------------------|:-----------------------
 snmp_host_name | string | - |v1,v2c,v3 | -  | no | SNMP host name give by the user, any string helpful for the user, this parameter is added as tag (SNMP_HOST_NAME) for metrics   
 snmp_host_address | string | - | v1,v2c,v3 | - | yes | IP address or host name with port number. This parameter is added as a tag (SNMP_HOST_ADDRESS) for metrics
 snmp_version | string | v1/v2c/v3 | v1,v2c,v3 | -  | yes | SNMP version
 community | string |  public/private | v1,v2c | - | yes | Community
 user_name | string | - |  v3 |  - | yes | User name
 security_level | string | noAuthNoPriv/authNoPriv/authPriv | v3 | - | yes | Security level
 auth_password | string | - | v3 | - | no |  Authentication protocol pass phrase
 auth_protocol | string | MD5/SHA | v3 | - | yes | Authentication protocol
 priv_password | string | - | v3 | - | no | Privacy protocol pass phrase
 priv_protocol | string | DES/AES| v3 | - | yes | Privacy protocol
 security_engine_id | string| - | v3 | - | no | Security engine ID
 context_engine_id | string | - | v3 | - | no | Context engine ID
 context_name | string | - | v3 | - | no | Context name 
 retries | uint | - | v1,v2c,v3 | 1 | no | Number of connection retries 
 timeout | int | -  | v1,v2c,v3 | 5 | no | SNMP request timeout in seconds

Exemplary task manifest (more examples in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/)):
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "120s"
    },
    "workflow": {
        "collect": {
            "metrics": {
				"/intel/snmp/*": {}
           },
            "config": {
				"/intel/snmp": {
                    "snmp_host_name": "snmp_host",
                    "snmp_host_address": "127.0.0.1",
					"user_name": "name",
					"security_level ": "authPriv",
					"auth_protocol": "MD5",
					"auth_passphrase": "password",
					"privacy_protocol": "DES",
					"privacy_passphrase": "password"
                }
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published.txt"
                    }
                }
            ]
        }
    }
}
```


### Examples
Example running snap-plugin-collector-snmp plugin and writing data to a file.

Create configuration file (setfile) for snmp plugin, see exemplary in [examples/setfiles/](https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/setfiles/).

Set path to configuration file as the field `setfile` in Global Config, see exemplary in [examples/configs/] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/configs/).

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in config.json ):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config config.json
```
In another terminal window:

Load snap-plugin-collector-snmp plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-snmp
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics for your system

```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-collector-snmp plugin (exemplary in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-snmp/blob/master/examples/tasks/)).

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out on other use cases, visit:
* [Snap Gitter channel](https://gitter.im/intelsdi-x/snap)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Katarzyna Zabrocka](https://github.com/katarzyna-z)

