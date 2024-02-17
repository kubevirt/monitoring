package prometheus

import "github.com/prometheus/client_golang/prometheus"

func no_local_name() {
	prometheus.Register(prometheus.NewGauge( // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory, using operator-observability packages.`
		prometheus.GaugeOpts{
			Name: "gauge_name",
			Help: "gauge_help",
		}))

	prometheus.MustRegister(prometheus.NewCounter( // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory, using operator-observability packages.`
		prometheus.CounterOpts{
			Name: "counter_name",
			Help: "counter_help",
		}))

	prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "hist_name",
		Help: "hist_help",
	})
}
