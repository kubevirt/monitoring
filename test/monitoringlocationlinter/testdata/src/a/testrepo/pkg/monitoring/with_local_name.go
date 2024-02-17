package metrics

import prom "github.com/prometheus/client_golang/prometheus"

func with_local_name() {
	prom.Register(prom.NewGauge(
		prom.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))
	prom.MustRegister(prom.NewCounter(
		prom.CounterOpts{
			Name: "counter_name",
			Help: "counter_help",
		}))
}
