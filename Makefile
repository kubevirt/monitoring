# Default to podman
CONTAINER_RUNTIME ?= podman

.PHONY: metricsdocs
metricsdocs: build-metricsdocs
	@[ "${CONFIG_FILE}" ] || ( echo "CONFIG_FILE is not set"; exit 1 )

	tools/metricsdocs/_out/metricsdocs \
		--config-file $(CONFIG_FILE)

.PHONY: build-metricsdocs
build-metricsdocs:
	cd ./tools/metricsdocs && go build -ldflags="-s -w" -o _out/metricsdocs .

.PHONY: promlinter-build
promlinter-build:
	${CONTAINER_RUNTIME} build -t ${IMG} test/metrics/prom-metrics-linter

.PHONY: promlinter-push
promlinter-push:
	${CONTAINER_RUNTIME} push ${IMG}
