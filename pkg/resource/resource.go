package resource

import "github.com/statsd/client-interface"

type Resource interface {
	Name() string
	Start(statsd.Client) error
	Stop() error
}
