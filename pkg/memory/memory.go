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

			m.client.Gauge("total", bytes(stat["MemTotal"]))
			m.client.Gauge("free", bytes(stat["MemFree"]))
			m.client.Gauge("active", bytes(stat["Active"]))
			m.client.Gauge("swap", bytes(stat["SwapTotal"]))
			m.client.Gauge("percent", percent(stat))

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

func percent(s linux.MemInfo) int {
	total := s["MemTotal"]
	used := total - s["MemFree"] - s["Buffers"] - s["Cached"]
	return int(float64(used) / float64(total) * 100)
}

// convert to bytes.
func bytes(n uint64) int {
	return int(n * 1000)
}
