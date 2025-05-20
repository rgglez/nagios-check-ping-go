# nagios-check-ping-go

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/nagios-check-ping-go/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/nagios-check-ping-go)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/nagios-check-ping-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/nagios-check-ping-go)](https://goreportcard.com/report/github.com/rgglez/nagios-check-ping-go)
[![GitHub release](https://img.shields.io/github/release/rgglez/nagios-check-ping-go.svg)](https://github.com/rgglez/nagios-check-ping-go/releases/)

*nagios-check-ping-go* is a plugin for [Nagios](https://www.nagios.org) written in [Go](https://go.dev/). It pings a host to check if it is alive. Currently this plugin uses [pro-bing](https://github.com/prometheus-community/pro-bing) and can be run in non-root mode.

## Usage

### Example

```bash
check_ping_go --host=www.example.com --warn=100,5% --crit=200,10% -count=5
```

This command pings www.example.com, reporting a warning if the ping times are above 100ms or 5% of packets are lost, and a critical condition if times excedes 200 or at least 10% are lost. It sends 5 packets.

### Command line parameters

* `--host` or `-H` string. Hostname or IP address to ping.
* `--warn` or `-w` string. Warning threshold in the form 'time,packets%'. Default: "150,5%".
* `--crit` or `-c` string. Critical threshold in the form 'time,packets%'. Default: "200,10%".
* `--count` or `-n` string. Number of packets to be send. Default: 5.

## Build and installation

### Dependencies

* [github.com/spf13/pflag](https://github.com/spf13/pflag)
* [github.com/xorpaul/go-nagios](https://github.com/xorpaul/go-nagios)
* [github.com/prometheus-community/pro-bing](github.com/prometheus-community/pro-bing)

### Get the code

```bash
$ git clone https://github.com/rgglez/nagios-check-ping-go.git
$ cd nagios-check-ping-go
```

### Build

To build the program, run:

```bash
$ make build
```

The executable will be created inside the ```dist``` directory.

### Installation

To install the binary to the default path (```/usr/local/nagios/libexec```), execute:

```bash
# make install
```

Or just copy the executable to your regular Nagios plugins directory.

### Permissions

While pro-bing was selected to minimize the requirements to run the plugin as root or with setuid, you still need to make some adjustments to your Linux system in order to run the program. 

As stated in its [documentation](https://github.com/prometheus-community/pro-bing/blob/main/README.md):

> This library attempts to send an "unprivileged" ping via UDP. On Linux, this must be enabled with the following sysctl command:

```bash
# sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
```

## License

Copyright 2024 Rodolfo González González.

[GPL v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). Please read the [LICENSE](LICENSE.md) file.

