package main

import "github.com/statsd/system/pkg/collector"
import "github.com/statsd/system/pkg/memory"
import "github.com/statsd/client-namespace"
import "github.com/statsd/system/pkg/disk"
import "github.com/statsd/system/pkg/cpu"
import . "github.com/tj/go-gracefully"
import "github.com/segmentio/go-log"
import "github.com/statsd/client"
import "github.com/tj/docopt"
import "time"
import "os"

const Version = "0.3.0"

const Usage = `
  Usage:
    system-stats
      [--statsd-address addr]
      [--memory-interval i]
      [--disk-interval i]
      [--cpu-interval i]
      [--extended]
      [--name name]
    system-stats -h | --help
    system-stats --version

  Options:
    --statsd-address addr   statsd address [default: :8125]
    --memory-interval i     memory reporting interval [default: 10s]
    --disk-interval i       disk reporting interval [default: 30s]
    --cpu-interval i        cpu reporting interval [default: 5s]
    --extended              output additional extended metrics
    --name name             node name defaulting to hostname [default: hostname]
    -h, --help              output help information
    -v, --version           output version
`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	log.Info("starting system %s", Version)

	client, err := statsd.Dial(args["--statsd-address"].(string))
	log.Check(err)

	extended := args["--extended"].(bool)

	name := args["--name"].(string)
	if "hostname" == name {
		host, err := os.Hostname()
		log.Check(err)
		name = host
	}

	c := collector.New(namespace.New(client, name))
	c.Add(memory.New(interval(args, "--memory-interval"), extended))
	c.Add(cpu.New(interval(args, "--cpu-interval"), extended))
	c.Add(disk.New(interval(args, "--disk-interval")))

	c.Start()
	Shutdown()
	c.Stop()
}

func interval(args map[string]interface{}, name string) time.Duration {
	d, err := time.ParseDuration(args[name].(string))
	log.Check(err)
	return d
}
