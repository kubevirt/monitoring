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
