package images

import (
	"github.com/prometheus/client_golang/prometheus"
)

var FunctionLatency = CreateExecutionTimeMetric("time spent of /image")

func CreateExecutionTimeMetric(help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "execution_latency_seconds",
		Help:    help,
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"step"})
}
