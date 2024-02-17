package controller

import prom "github.com/prometheus/client_golang/prometheus"

func with_local_name() {
	prom.Register(prom.NewGauge( // want `monitoring-location-linter: Prometheus metrics should be registered only under pkg/monitoring.`
		prom.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))
	prom.MustRegister(prom.NewCounter( // want `monitoring-location-linter: Prometheus metrics should be registered only under pkg/monitoring.`
		prom.CounterOpts{
			Name: "counter_name",
			Help: "counter_help",
		}))
}
