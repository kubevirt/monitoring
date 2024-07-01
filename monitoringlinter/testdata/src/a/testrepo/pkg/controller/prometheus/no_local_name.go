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
