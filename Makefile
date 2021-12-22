.PHONY: metricsdocs build-metricsdocs

metricsdocs: build-metricsdocs
	@[ "${GITHUB_TOKEN_FILE}" ] || ( echo "GITHUB_TOKEN_FILE is not set"; exit 1 )
	@[ "${CONFIG_FILE}" ] || ( echo "CONFIG_FILE is not set"; exit 1 )

	tools/metricsdocs/_out/metricsdocs \
		--github-token-file $(GITHUB_TOKEN_FILE) \
		--config-file $(CONFIG_FILE)

build-metricsdocs:
	cd ./tools/metricsdocs && go build -ldflags="-s -w" -o _out/metricsdocs .
