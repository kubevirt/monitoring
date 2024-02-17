package operatorobservability

import (
	"github.com/machadovilaca/operator-observability/pkg/operatormetrics"
	"github.com/machadovilaca/operator-observability/pkg/operatorrules"
)

func no_local_name() {
	operatormetrics.RegisterMetrics()
	operatorrules.RegisterAlerts()
	operatorrules.RegisterRecordingRules()
}
