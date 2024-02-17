package prometheus

import prom "github.com/prometheus/client_golang/prometheus"

func with_local_name() {
	prom.Register(prom.NewGauge( // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory, using operator-observability packages.`
		prom.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))

	prom.MustRegister(prom.NewCounter( // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory, using operator-observability packages.`
		prom.CounterOpts{
			Name: "counter_name",
			Help: "counter_help",
		}))

	prom.NewHistogram(prom.HistogramOpts{
		Name: "hist_name",
		Help: "hist_help",
	})
}
