package controller

import "github.com/prometheus/client_golang/prometheus"

func no_local_name() {
	prometheus.Register(prometheus.NewGauge( // want `monitoring-location-linter: Prometheus metrics should be registered only under pkg/monitoring.`
		prometheus.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))
	prometheus.MustRegister(prometheus.NewCounter( // want `monitoring-location-linter: Prometheus metrics should be registered only under pkg/monitoring.`
		prometheus.CounterOpts{
			Name: "counter_name",
			Help: "counter_help",
		}))
}
