
# System

 System statistics collector for statsd written in Go.

## Usage

```

  Usage:
    system
      [--statsd-address addr]
      [--memory-interval i]
      [--disk-interval i]
      [--cpu-interval i]
      [--extended]
      [--name name]
    system -h | --help
    system --version

  Options:
    --statsd-address addr   statsd address [default: :8125]
    --memory-interval i     memory reporting interval [default: 10s]
    --disk-interval i       disk reporting interval [default: 1m]
    --cpu-interval i        cpu reporting interval [default: 2s]
    --name name             node name defaulting to hostname [default: hostname]
    --extended              output additional extended metrics
    -h, --help              output help information
    -v, --version           output version

````

## Installation

 Via go-get:

```
$ go get github.com/statsd/system
```

 Via binaries:

Coming soon!

## Metrics

 Metrics are prefixed with the hostname (or `--name`), and
 namespaced by the resource, for example:

```
api-2.cpu.blocked:7|c
api-2.cpu.running:4|c
api-2.cpu.interrupts:19695796035|c
api-2.cpu.percent:26|g
api-2.cpu.switches:25195265352|c
api-2.cpu.processes:20027|c
api-2.cpu.blocked:7|c
api-2.cpu.running:4|c
api-2.cpu.interrupts:19695796035|c
api-2.cpu.percent:26|g
...
```

 Depending on the statd implementation that you use this
 may result in different outputs. For example with the
 [armon/statsite](https://github.com/armon/statsite) implementation
 this would result in `gauges.api-2.memory.free`.

 The `--extended` flag enables extended metrics per resource
 and are listed as __extended__ below.

### CPU

- `cpu.percent` gauge
- `cpu.switches` counter __extended__
- `cpu.interrupts` counter __extended__
- `cpu.blocked` counter __extended__

### Memory

 Memory values are represented in bytes.

- `memory.percent` gauge
- `memory.used` gauge
- `memory.active` gauge __extended__
- `memory.total` gauge __extended__
- `memory.free` gauge __extended__
- `memory.swap.percent` gauge
- `memory.swap.total` gauge __extended__
- `memory.swap.free` gauge __extended__

### Disk

 Disk values are represented in bytes. `<volume>` is the
 path the fs is mounted on (/, /data, etc).

- `disk.<volume>.percent` gauge
- `disk.<volume>.free` gauge
- `disk.<volume>.used` gauge

### IO

  Coming soon!

## Daemonization

 system(1) doesn't support running as a daemon natively, you'll
 want to use upstart or similar for this. Add the following example
 upstart script to /etc/init/system.conf:

```
respawn

start on runlevel [2345]
stop on runlevel [016]

exec system --statsd-address 10.0.0.214:5000
```

 Then run `sudo start system` and you're good to go!

## Debugging

Run with `DEBUG=stats` to view the [go-debug](http://github.com/visionmedia/go-debug) output:

```
2014-08-13 22:04:36 INFO - cpu: reporting
22:04:36.098 2s     2s     statsd - vagrant-ubuntu-precise-64.cpu.switches:20384|c
22:04:36.098 4us    4us    statsd - vagrant-ubuntu-precise-64.cpu.processes:0|c
22:04:36.098 3us    3us    statsd - vagrant-ubuntu-precise-64.cpu.running:0|c
22:04:36.098 3us    3us    statsd - vagrant-ubuntu-precise-64.cpu.interrupts:656|c
22:04:36.098 3us    3us    statsd - vagrant-ubuntu-precise-64.cpu.percent:100|g
2014-08-13 22:04:38 INFO - cpu: reporting
22:04:38.098 2s     2s     statsd - vagrant-ubuntu-precise-64.cpu.switches:24074|c
22:04:38.098 23us   13us   statsd - vagrant-ubuntu-precise-64.cpu.processes:0|c
22:04:38.098 15us   8us    statsd - vagrant-ubuntu-precise-64.cpu.running:1|c
22:04:38.098 12us   7us    statsd - vagrant-ubuntu-precise-64.cpu.interrupts:638|c
22:04:38.099 11us   7us    statsd - vagrant-ubuntu-precise-64.cpu.percent:100|g
```

# License

 MIT
