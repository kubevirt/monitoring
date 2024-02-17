package operatorobservability

import (
	metrics "github.com/machadovilaca/operator-observability/pkg/operatormetrics"
	rules "github.com/machadovilaca/operator-observability/pkg/operatorrules"
)

func with_local_name() {
	metrics.RegisterMetrics() // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory.`
	metrics.ListMetrics()

	rules.RegisterAlerts() // want `monitoring-linter: alerts and recording rules should be registered only within pkg/monitoring directory.`
	rules.ListAlerts()

	rules.RegisterRecordingRules() // want `monitoring-linter: alerts and recording rules should be registered only within pkg/monitoring directory.`
	rules.ListRecordingRules()
}
