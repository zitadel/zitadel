go_bin := "$(shell go env GOPATH)/bin"
gen_authopt_path := "$(go_bin)/protoc-gen-authoption"
gen_zitadel_path := "$(go_bin)/protoc-gen-zitadel"

now := $(shell date '+%Y-%m-%dT%T%z' | sed -E 's/.([0-9]{2})([0-9]{2})$$/-\1:\2/')
VERSION ?= development-$(now)
COMMIT_SHA ?= $(shell git rev-parse HEAD)
ZITADEL_IMAGE ?= zitadel:local

GOCOVERDIR = tmp/coverage
ZITADEL_MASTERKEY ?= MasterkeyNeedsToHave32Characters

export GOCOVERDIR ZITADEL_MASTERKEY

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

export GOOS GOARCH

.PHONY: clean
clean:
	$(RM) -r .artifacts/grpc
	$(RM) $(gen_authopt_path)
	$(RM) $(gen_zitadel_path)
	$(RM) -r tmp/

.PHONY: api_integration_db_up
api_integration_db_up:
	docker compose -f internal/integration/config/docker-compose.yaml up --pull always --wait cache postgres

.PHONY: api_integration_db_down
api_integration_db_down:
	docker compose -f internal/integration/config/docker-compose.yaml down -v

.PHONY: api_integration_setup
api_integration_setup:
	go build -cover -race -tags integration -o zitadel.test main.go
	mkdir -p $${GOCOVERDIR}
	GORACE="halt_on_error=1" ./zitadel.test init --config internal/integration/config/zitadel.yaml --config internal/integration/config/postgres.yaml
	GORACE="halt_on_error=1" ./zitadel.test setup --masterkeyFromEnv --init-projections --config internal/integration/config/zitadel.yaml --config internal/integration/config/postgres.yaml --steps internal/integration/config/steps.yaml

.PHONY: api_integration_server_start
api_integration_server_start: api_integration_setup
	GORACE="log_path=tmp/race.log" \
	./zitadel.test start --masterkeyFromEnv --config internal/integration/config/zitadel.yaml --config internal/integration/config/postgres.yaml \
	  > tmp/zitadel.log 2>&1 \
	  & printf $$! > tmp/zitadel.pid

.PHONY: api_integration_test_packages
api_integration_test_packages:
	go test -race -count 1 -tags integration -timeout 60m -parallel 1 $$(go list -tags integration ./... | grep -e "integration_test" -e "events_testing")

.PHONY: api_integration_server_stop
api_integration_server_stop:
	pid=$$(cat tmp/zitadel.pid); \
	$(RM) tmp/zitadel.pid; \
	kill $$pid; \
	if [ -s tmp/race.log.$$pid ]; then \
		cat tmp/race.log.$$pid; \
		exit 66; \
	fi

.PHONY: api_integration_reports
api_integration_reports:
	go tool covdata textfmt -i=tmp/coverage -pkg=github.com/zitadel/zitadel/internal/...,github.com/zitadel/zitadel/cmd/...,github.com/zitadel/zitadel/backend/... -o profile.cov

.PHONY: api_integration_test
api_integration_test: api_integration_server_start api_integration_test_packages api_integration_server_stop api_integration_reports
