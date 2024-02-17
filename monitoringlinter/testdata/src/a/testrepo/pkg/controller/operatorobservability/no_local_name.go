package operatorobservability

import (
	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
	"github.com/machadovilaca/operator-observability/pkg/operatorrules"
)

func no_local_name() {
	operatormetrics.RegisterMetrics() // want `monitoring-linter: metrics should be registered only within pkg/monitoring directory.`
	operatormetrics.ListMetrics()

	operatorrules.RegisterAlerts() // want `monitoring-linter: alerts and recording rules should be registered only within pkg/monitoring directory.`
	operatorrules.ListAlerts()

	operatorrules.RegisterRecordingRules() // want `monitoring-linter: alerts and recording rules should be registered only within pkg/monitoring directory.`
	operatorrules.ListRecordingRules()
}
