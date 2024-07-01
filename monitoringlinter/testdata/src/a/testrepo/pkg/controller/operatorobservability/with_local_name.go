/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2024 Red Hat, Inc.
 *
 */

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
