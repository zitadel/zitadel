go_bin := "$$(go env GOPATH)/bin"
gen_authopt_path := "$(go_bin)/protoc-gen-authoption"
gen_zitadel_path := "$(go_bin)/protoc-gen-zitadel"

now := $(shell date '+%Y-%m-%dT%T%z' | sed -E 's/.([0-9]{2})([0-9]{2})$$/-\1:\2/')
VERSION ?= development-$(now)
COMMIT_SHA ?= $(shell git rev-parse HEAD)
ZITADEL_IMAGE ?= zitadel:local

GOCOVERDIR = tmp/coverage
INTEGRATION_DB_FLAVOR ?= postgres
ZITADEL_MASTERKEY ?= MasterkeyNeedsToHave32Characters

export GOCOVERDIR INTEGRATION_DB_FLAVOR ZITADEL_MASTERKEY

.PHONY: compile
compile: core_build console_build compile_pipeline

.PHONY: docker_image
docker_image: compile
	DOCKER_BUILDKIT=1 docker build -f build/Dockerfile -t $(ZITADEL_IMAGE) .

.PHONY: compile_pipeline
compile_pipeline: console_move
	CGO_ENABLED=0 go build -o zitadel -v -ldflags="-s -w -X 'github.com/zitadel/zitadel/cmd/build.commit=$(COMMIT_SHA)' -X 'github.com/zitadel/zitadel/cmd/build.date=$(now)' -X 'github.com/zitadel/zitadel/cmd/build.version=$(VERSION)' "
	chmod +x zitadel

.PHONY: core_dependencies
core_dependencies:
	go mod download

.PHONY: core_static
core_static:
	go install github.com/rakyll/statik@v0.1.7
	go generate internal/api/ui/login/static/resources/generate.go
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go

.PHONY: core_generate_all
core_generate_all:
	go install github.com/dmarkham/enumer@v1.5.10 		# https://pkg.go.dev/github.com/dmarkham/enumer?tab=versions
	go install github.com/rakyll/statik@v0.1.7			# https://pkg.go.dev/github.com/rakyll/statik?tab=versions
	go install go.uber.org/mock/mockgen@v0.4.0			# https://pkg.go.dev/go.uber.org/mock/mockgen?tab=versions
	go install golang.org/x/tools/cmd/stringer@v0.22.0	# https://pkg.go.dev/golang.org/x/tools/cmd/stringer?tab=versions
	go generate ./...

.PHONY: core_assets
core_assets:
	mkdir -p docs/apis/assets
	go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

.PHONY: core_api_generator
core_api_generator:
ifeq (,$(wildcard $(gen_authopt_path)))
	go install internal/protoc/protoc-gen-authoption/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_authopt_path)
endif
ifeq (,$(wildcard $(gen_zitadel_path)))
	go install internal/protoc/protoc-gen-zitadel/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_zitadel_path)
endif

.PHONY: core_grpc_dependencies
core_grpc_dependencies:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1 						# https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go?tab=versions
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 						# https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0	# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0 		# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2?tab=versions
	go install github.com/envoyproxy/protoc-gen-validate@v1.1.0								# https://pkg.go.dev/github.com/envoyproxy/protoc-gen-validate?tab=versions
	go install github.com/bufbuild/buf/cmd/buf@v1.45.0										# https://pkg.go.dev/github.com/bufbuild/buf/cmd/buf?tab=versions

.PHONY: core_api
core_api: core_api_generator core_grpc_dependencies
	buf generate
	mkdir -p pkg/grpc
	cp -r .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc/
	mkdir -p openapi/v2/zitadel
	cp -r .artifacts/grpc/zitadel/ openapi/v2/zitadel

.PHONY: core_build
core_build: core_dependencies core_api core_static core_assets

.PHONY: console_move
console_move:
	cp -r console/dist/console/* internal/api/ui/console/static

.PHONY: console_dependencies
console_dependencies:
	cd console && \
	yarn install --immutable

.PHONY: console_client
console_client:
	cd console && \
	yarn generate

.PHONY: console_build
console_build: console_dependencies console_client
	cd console && \
	yarn build

.PHONY: clean
clean:
	$(RM) -r .artifacts/grpc
	$(RM) $(gen_authopt_path)
	$(RM) $(gen_zitadel_path)
	$(RM) -r tmp/

.PHONY: core_unit_test
core_unit_test:
	go test -race -coverprofile=profile.cov -coverpkg=./internal/...  ./...

.PHONY: core_integration_db_up
core_integration_db_up:
	docker compose -f internal/integration/config/docker-compose.yaml up --pull always --wait $${INTEGRATION_DB_FLAVOR} cache

.PHONY: core_integration_db_down
core_integration_db_down:
	docker compose -f internal/integration/config/docker-compose.yaml down

.PHONY: core_integration_setup
core_integration_setup:
	go build -cover -race -tags integration -o zitadel.test main.go
	mkdir -p $${GOCOVERDIR}
	GORACE="halt_on_error=1" ./zitadel.test init --config internal/integration/config/zitadel.yaml --config internal/integration/config/${INTEGRATION_DB_FLAVOR}.yaml
	GORACE="halt_on_error=1" ./zitadel.test setup --masterkeyFromEnv --init-projections --config internal/integration/config/zitadel.yaml --config internal/integration/config/${INTEGRATION_DB_FLAVOR}.yaml --steps internal/integration/config/steps.yaml

.PHONY: core_integration_server_start
core_integration_server_start: core_integration_setup
	GORACE="log_path=tmp/race.log" \
	./zitadel.test start --masterkeyFromEnv --config internal/integration/config/zitadel.yaml --config internal/integration/config/${INTEGRATION_DB_FLAVOR}.yaml \
	  > tmp/zitadel.log 2>&1 \
	  & printf $$! > tmp/zitadel.pid

.PHONY: core_integration_test_packages
core_integration_test_packages:
	go test -race -count 1 -tags integration -timeout 30m $$(go list -tags integration ./... | grep "integration_test")

.PHONY: core_integration_server_stop
core_integration_server_stop:
	pid=$$(cat tmp/zitadel.pid); \
	$(RM) tmp/zitadel.pid; \
	kill $$pid; \
	if [ -s tmp/race.log.$$pid ]; then \
		cat tmp/race.log.$$pid; \
		exit 66; \
	fi

.PHONY: core_integration_reports
core_integration_reports:
	go tool covdata textfmt -i=tmp/coverage -pkg=github.com/zitadel/zitadel/internal/...,github.com/zitadel/zitadel/cmd/... -o profile.cov

.PHONY: core_integration_test
core_integration_test: core_integration_server_start core_integration_test_packages core_integration_server_stop core_integration_reports

.PHONY: console_lint
console_lint:
	cd console && \
	yarn lint

.PHONY: core_lint
core_lint:
	golangci-lint run \
		--timeout 10m \
		--config ./.golangci.yaml \
		--out-format=github-actions \
		--concurrency=$$(getconf _NPROCESSORS_ONLN)
