package prometric

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewMetric() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests"},
		[]string{"method"})
}
