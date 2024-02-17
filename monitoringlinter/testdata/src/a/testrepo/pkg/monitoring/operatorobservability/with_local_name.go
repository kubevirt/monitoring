package operatorobservability

import (
	metrics "github.com/machadovilaca/operator-observability/pkg/operatormetrics"
	rules "github.com/machadovilaca/operator-observability/pkg/operatorrules"
)

func with_local_name() {
	metrics.RegisterMetrics()
	rules.RegisterAlerts()
	rules.RegisterRecordingRules()
}
