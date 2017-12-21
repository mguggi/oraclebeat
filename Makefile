BEAT_NAME=oraclebeat
BEAT_PATH=github.com/mguggi/oraclebeat
BEAT_URL=https://${BEAT_PATH}
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell govendor list -no-status +local)
PREFIX?=.
NOTICE_FILE=NOTICE

# Overwrite docker compose project name
DOCKER_COMPOSE_PROJECT_NAME?=${BEAT_NAME}_

# Load environment file and export variable(s)
include module/oracle/_meta/env
export $(shell sed -e '/^\#/d' -e 's/=.*//' module/oracle/_meta/env)

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make create-metricset
	make collect

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:
	if [ -f modules.d/oracle.yml.disabled ]; then mv modules.d/oracle.yml.disabled modules.d/oracle.yml; fi

# Test only local
.PHONY: local-tests
local-tests: ## @testing Runs the unit tests with coverage.
    ${MAKE} prepare-tests
	go test -i ${GOPACKAGES}
	$(COVERAGE_TOOL) $(RACE) -coverprofile=${COVERAGE_DIR}/unit.cov  ${GOPACKAGES}
