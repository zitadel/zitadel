go_bin := "$$(go env GOPATH)/bin"
gen_authopt_path := "$(go_bin)/protoc-gen-authoption"
gen_zitadel_path := "$(go_bin)/protoc-gen-zitadel"

now := $(shell date '+%Y-%m-%dT%T%z' | sed -E 's/.([0-9]{2})([0-9]{2})$$/-\1:\2/')
VERSION ?= development-$(now)
COMMIT_SHA ?= $(shell git rev-parse HEAD)
ZITADEL_IMAGE ?= zitadel:local

GOCOVERDIR = tmp/coverage
ZITADEL_MASTERKEY ?= MasterkeyNeedsToHave32Characters

export GOCOVERDIR ZITADEL_MASTERKEY

GOOS?=$(go env GOOS)
GOARCH?=$(go env GOARCH)

export GOOS GOARCH

.PHONY: compile
compile: api_build console_build compile_pipeline

.PHONY: docker_image
docker_image:
	@if [ ! -f .artifacts/${GOOS}/${GOARCH}/zitadel ]; then \
		echo "Compiling zitadel binary"; \
		$(MAKE) compile; \
	else \
		echo "Reusing precompiled zitadel binary"; \
	fi
	DOCKER_BUILDKIT=1 docker build -f apps/api/Dockerfile -t $(ZITADEL_IMAGE) .

.PHONY: compile_pipeline
compile_pipeline: console_move
	CGO_ENABLED=0 go build -o .artifacts/${GOOS}/${GOARCH}/zitadel -v -ldflags="-s -w -X 'github.com/zitadel/zitadel/cmd/build.commit=$(COMMIT_SHA)' -X 'github.com/zitadel/zitadel/cmd/build.date=$(now)' -X 'github.com/zitadel/zitadel/cmd/build.version=$(VERSION)' "
	chmod +x .artifacts/${GOOS}/${GOARCH}/zitadel

.PHONY: api_dependencies
api_dependencies:
	go mod download

.PHONY: api_static
api_static:
	go install github.com/rakyll/statik@v0.1.7
	go generate internal/api/ui/login/static/resources/generate.go
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go

.PHONY: api_generate_all
api_generate_all:
	go install github.com/dmarkham/enumer@v1.5.11 		# https://pkg.go.dev/github.com/dmarkham/enumer?tab=versions
	go install github.com/rakyll/statik@v0.1.7			# https://pkg.go.dev/github.com/rakyll/statik?tab=versions
	go install go.uber.org/mock/mockgen@v0.4.0			# https://pkg.go.dev/go.uber.org/mock/mockgen?tab=versions
	go install golang.org/x/tools/cmd/stringer@v0.36.0	# https://pkg.go.dev/golang.org/x/tools/cmd/stringer?tab=versions
	go generate ./...

.PHONY: api_assets
api_assets:
	mkdir -p docs/apis/assets
	go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

.PHONY: api_stubs_generator
api_stubs_generator:
ifeq (,$(wildcard $(gen_authopt_path)))
	go install internal/protoc/protoc-gen-authoption/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_authopt_path)
endif
ifeq (,$(wildcard $(gen_zitadel_path)))
	go install internal/protoc/protoc-gen-zitadel/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_zitadel_path)
endif

.PHONY: api_grpc_dependencies
api_grpc_dependencies:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1 						# https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go?tab=versions
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 						# https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0	# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway?tab=versions
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0 		# https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2?tab=versions
	go install github.com/envoyproxy/protoc-gen-validate@v1.1.0								# https://pkg.go.dev/github.com/envoyproxy/protoc-gen-validate?tab=versions
	go install github.com/bufbuild/buf/cmd/buf@v1.45.0										# https://pkg.go.dev/github.com/bufbuild/buf/cmd/buf?tab=versions
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.18.1						# https://pkg.go.dev/connectrpc.com/connect/cmd/protoc-gen-connect-go?tab=versions

.PHONY: api_stubs
api_stubs: api_stubs_generator api_grpc_dependencies
	buf generate
	mkdir -p pkg/grpc
	cp -r .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/** pkg/grpc/
	mkdir -p openapi/v2/zitadel
	cp -r .artifacts/grpc/zitadel/ openapi/v2/zitadel

.PHONY: api_build
api_build: api_dependencies api_stubs api_static api_assets

.PHONY: console_move
console_move:
	cp -r apps/console/dist/console/* internal/api/ui/console/static

.PHONY: console_dependencies
console_dependencies:
	npx pnpm install --frozen-lockfile --filter=console...

.PHONY: console_build
console_build: console_dependencies
	nx run @zitadel/console:build

.PHONY: clean
clean:
	$(RM) -r .artifacts/grpc
	$(RM) $(gen_authopt_path)
	$(RM) $(gen_zitadel_path)
	$(RM) -r tmp/

.PHONY: api_unit_test
api_unit_test:
	go test -race -coverprofile=profile.cov -coverpkg=./internal/...  ./...

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

.PHONY: console_lint
console_lint:
	nx run @zitadel/console:lint

.PHONY: api_lint
api_lint:
	golangci-lint run \
		--timeout 10m \
		--config ./.golangci.yaml \
		--out-format=github-actions \
		--concurrency=$$(getconf _NPROCESSORS_ONLN)
