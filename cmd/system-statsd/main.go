package main

import (
	"os"
	"time"

	"github.com/ques0942/system/pkg/collector"
	"github.com/ques0942/system/pkg/cpu"
	"github.com/ques0942/system/pkg/disk"
	"github.com/ques0942/system/pkg/memory"
	"github.com/segmentio/go-log"
	"github.com/statsd/client"
	"github.com/statsd/client-namespace"
	"github.com/tj/docopt"
	. "github.com/tj/go-gracefully"
)

const Version = "0.2.0"

const Usage = `
  Usage:
    system-statsd
      [--statsd-address addr]
      [--memory-interval i]
      [--disk-interval i]
      [--cpu-interval i]
      [--extended]
      [--name name]
    system-statsd -h | --help
    system-statsd --version

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
