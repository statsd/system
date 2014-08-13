//
// Memory resource.
//
// This collector reports on the following meminfo metrics:
//
//  - "percent" (gauge)
//  - "active" (gauge)
//  - "total" (gauge)
//  - "free" (gauge)
//  - "swap" (gauge)
//
package memory

import "github.com/statsd/client-interface"
import "github.com/c9s/goprocinfo/linux"
import "github.com/segmentio/go-log"
import "time"

// Memory resource.
type Memory struct {
	Path     string
	Interval time.Duration
	client   statsd.Client
	exit     chan struct{}
}

// New memory resource.
func New(interval time.Duration) *Memory {
	return &Memory{
		Path:     "/proc/meminfo",
		exit:     make(chan struct{}),
		Interval: interval,
	}
}

// Name of the resource.
func (m *Memory) Name() string {
	return "memory"
}

// Start resource collection.
func (m *Memory) Start(client statsd.Client) error {
	m.client = client
	go m.Report()
	return nil
}

// Report resource.
func (m *Memory) Report() {
	for {
		select {
		case <-time.Tick(m.Interval):
			log.Info("memory: reporting")

			stat, err := linux.ReadMemInfo(m.Path)

			if err != nil {
				log.Error("memory: %s", err)
				continue
			}

			m.report(stat, "MemTotal", "total")
			m.report(stat, "MemFree", "free")
			m.report(stat, "Active", "active")
			m.report(stat, "SwapTotal", "swap")
			m.reportPercent(stat)

		case <-m.exit:
			log.Info("mem: exiting")
			return
		}
	}
}

// Stop resource collection.
func (m *Memory) Stop() error {
	close(m.exit)
	return nil
}

// report percentage.
func (m *Memory) reportPercent(stat linux.MemInfo) {
	if total, ok := stat["MemTotal"]; ok {
		if free, ok := stat["MemFree"]; ok {
			m.client.Gauge("percent", int(float64(total-free)/float64(total)*100))
		}
	}
}

// report the given `metric` as `name`.
func (m *Memory) report(stat linux.MemInfo, metric, name string) {
	if val, ok := stat[metric]; ok {
		m.client.Gauge(name, int(val*1000))
	}
}
