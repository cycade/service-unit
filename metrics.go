package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func RegisterMetricsFunc(c prometheus.HistogramVec) {
	err := prometheus.Register(c)
	if err != nil {
		fmt.Println(err)
	}
}
