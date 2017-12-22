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

# Test oracle module
.PHONY: test-oraclebeat
test-oraclebeat: python-env
	. ${PYTHON_ENV}/bin/activate && docker-compose build --force-rm oracle
	docker-compose run -d -p 1521:1521 oracle
	sleep 30 ## wait for oracle database service
	go test -tags=integration ${BEAT_PATH}/module/oracle/... -v
	. ${PYTHON_ENV}/bin/activate && INTEGRATION_TESTS=1 nosetests tests/system/test_oracle.py
	mkdir ${COVERAGE_DIR}
	$(COVERAGE_TOOL) $(RACE) -coverprofile=${COVERAGE_DIR}/unit.cov ${GOPACKAGES}
