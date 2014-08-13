//
// Disk resource.
//
// This collector reports on the following stat metrics:
//
//  - "disk.percent" (gauge)
//  - "disk.free" (gauge)
//  - "disk.used" (gauge)
//
package disk

import "github.com/statsd/client-interface"
import "github.com/c9s/goprocinfo/linux"
import "github.com/segmentio/go-log"
import "time"

// Disk resource.
type Disk struct {
	Path     string
	Interval time.Duration
	client   statsd.Client
	exit     chan struct{}
}

// New disk resource.
func New(interval time.Duration) *Disk {
	return &Disk{
		Path:     "/",
		Interval: interval,
		exit:     make(chan struct{}),
	}
}

// Name of resource.
func (d *Disk) Name() string {
	return "disk"
}

// Start resource collection.
func (d *Disk) Start(client statsd.Client) error {
	d.client = client
	go d.Report()
	return nil
}

// Report resources.
func (d *Disk) Report() {
	for {
		select {
		case <-time.Tick(d.Interval):
			log.Info("disk: reporting")
			stat, err := linux.ReadDisk(d.Path)

			if err != nil {
				log.Error("disk: %s", err)
				continue
			}

			d.client.Gauge("percent", int(percent(stat.Used, stat.All)))
			d.client.Gauge("free", int(stat.Free))
			d.client.Gauge("used", int(stat.Used))

		case <-d.exit:
			log.Info("disk: exiting")
			return
		}
	}
}

// Stop resource collection.
func (d *Disk) Stop() error {
	println("stopping disk")
	return nil
}

// calculate percentage.
func percent(a, b uint64) uint64 {
	return uint64(float64(a) / float64(b) * 100)
}
