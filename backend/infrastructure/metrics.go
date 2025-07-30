package infrastructure

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	CustomCounter prometheus.Counter
}

var Metrics = metrics{
	CustomCounter: promauto.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_counter",
		Help: "My Custom counter",
	}),
}
