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