.PHONY: metricsdocs build-metricsdocs

metricsdocs: build-metricsdocs
	@[ "${CONFIG_FILE}" ] || ( echo "CONFIG_FILE is not set"; exit 1 )

	tools/metricsdocs/_out/metricsdocs \
		--config-file $(CONFIG_FILE)

build-metricsdocs:
	cd ./tools/metricsdocs && go build -ldflags="-s -w" -o _out/metricsdocs .
