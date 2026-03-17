# ZITADEL Go Build Makefile
# This Makefile encapsulates Go-specific build concerns.
# NX delegates to these targets for Go builds, but they can also be run directly.

# ─── Configuration ───────────────────────────────────────────────────────────

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BIN_DIR = .artifacts/bin
TOOL_DIR = $(BIN_DIR)/$(GOOS)/$(GOARCH)

VERSION ?= local
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date "+%Y-%m-%dT%T%z" | sed -E 's/.([0-9]{2})([0-9]{2})$$/-\1:\2/')

LDFLAGS = -s -w \
	-X github.com/zitadel/zitadel/cmd/build.version=$(VERSION) \
	-X github.com/zitadel/zitadel/cmd/build.commit=$(COMMIT) \
	-X github.com/zitadel/zitadel/cmd/build.date=$(DATE)

# ─── Generate ────────────────────────────────────────────────────────────────

.PHONY: generate-install
generate-install:
	@echo "==> Installing Go tools..."
	@mkdir -p $(TOOL_DIR)
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/daixiang0/gci@v0.11.2
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/dmarkham/enumer@v1.5.11
	GOBIN=$(PWD)/$(TOOL_DIR) go install go.uber.org/mock/mockgen@v0.4.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install golang.org/x/tools/cmd/stringer@v0.36.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/rakyll/statik@v0.1.7
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/bufbuild/buf/cmd/buf@v1.45.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1
	GOBIN=$(PWD)/$(TOOL_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(PWD)/$(TOOL_DIR) go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.18.1
	GOBIN=$(PWD)/$(TOOL_DIR) go install ./internal/protoc/protoc-gen-authoption
	GOBIN=$(PWD)/$(TOOL_DIR) go install ./internal/protoc/protoc-gen-zitadel

.PHONY: generate-stubs
generate-stubs: generate-install
	@echo "==> Generating gRPC and OpenAPI stubs..."
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" buf generate
	@mkdir -p pkg/grpc openapi/v2/zitadel
	cp -r .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/** pkg/grpc/
	cp -r .artifacts/grpc/zitadel/ openapi/v2/zitadel

.PHONY: generate-statik
generate-statik: generate-install
	@echo "==> Generating statik files..."
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" go generate internal/api/ui/login/static/resources/generate.go
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" go generate internal/api/ui/login/statik/generate.go
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" go generate internal/notification/statik/generate.go
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" go generate internal/statik/generate.go

.PHONY: generate-assets
generate-assets: generate-install
	@echo "==> Generating asset routes and documentation..."
	@mkdir -p apps/docs/content/apis/assets
	go run internal/api/assets/generator/asset_generator.go \
		-directory=internal/api/assets/generator/ \
		-assets=apps/docs/content/apis/assets/assets.mdx

.PHONY: generate
generate: generate-stubs generate-statik generate-assets
	@echo "==> All generation complete."

# ─── Build ───────────────────────────────────────────────────────────────────

.PHONY: embed-console
embed-console:
	@echo "==> Copying Console dist to API embed directory..."
	find internal/api/ui/console/static -mindepth 1 ! -name 'gitkeep' -delete
	cp -r console/dist/console/* internal/api/ui/console/static/

.PHONY: build
build:
	@echo "==> Building zitadel for $(GOOS)/$(GOARCH)..."
	CGO_ENABLED=0 go build \
		-o $(BIN_DIR)/$(GOOS)/$(GOARCH)/zitadel.local \
		-ldflags="$(LDFLAGS)" \
		.

.PHONY: build-cli
build-cli:
	@echo "==> Building zitadel-cli for $(GOOS)/$(GOARCH)..."
	CGO_ENABLED=0 go build \
		-o $(BIN_DIR)/$(GOOS)/$(GOARCH)/zitadel-cli.local \
		-ldflags="$(LDFLAGS)" \
		./backend/main.go

# ─── Release (Multi-Platform) ────────────────────────────────────────────────
# Multi-platform cross-compilation is handled by GoReleaser.
# Use 'goreleaser build --snapshot --clean' for local testing.
# Use 'goreleaser release --clean' for the full release pipeline.

# ─── Lint ────────────────────────────────────────────────────────────────────

.PHONY: lint-install
lint-install:
	@echo "==> Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $(TOOL_DIR) v2.11.3

.PHONY: lint
lint: lint-install
	@echo "==> Linting Go code..."
	PATH="$(PWD)/$(TOOL_DIR):$$PATH" golangci-lint run --timeout 15m --config ./.golangci.yaml --verbose

# ─── Test ────────────────────────────────────────────────────────────────────

.PHONY: test-unit
test-unit:
	@echo "==> Running unit tests with coverage (no race)..."
	go test -coverprofile=profile.api.test-unit.cov \
		-coverpkg=./internal/...,./backend/... \
		./...

.PHONY: test-unit-race
test-unit-race:
	@echo "==> Running unit tests with race detector (no coverage)..."
	go test -race ./...

.PHONY: test-integration
test-integration:
	@echo "==> Running integration tests (testcontainers)..."
	go test -count 1 -tags integration -timeout 60m -v ./tests/integration/...

# ─── Clean ───────────────────────────────────────────────────────────────────

.PHONY: clean
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BIN_DIR) .artifacts/pack
	rm -f profile.api.test-unit.cov profile.api.test-integration.cov

# ─── Help ────────────────────────────────────────────────────────────────────

.PHONY: help
help:
	@echo "ZITADEL Go Build Targets:"
	@echo ""
	@echo "  make generate          Generate all code (proto stubs, statik, assets)"
	@echo "  make build             Build the API binary for the current platform"
	@echo "  make build-cli         Build the CLI binary for the current platform"
	@echo "  make lint              Lint Go code with golangci-lint"
	@echo "  make test-unit         Run unit tests with coverage (no race)"
	@echo "  make test-unit-race    Run unit tests with race detector (no coverage)"
	@echo "  make test-integration  Run integration tests (testcontainers)"
	@echo "  make clean             Remove build artifacts"
	@echo ""
	@echo "Release (run inside devcontainer):"
	@echo "  goreleaser build --snapshot --clean   Local multi-platform build"
	@echo "  goreleaser release --clean            Full release pipeline"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)     Version string for ldflags"
	@echo "  GOOS=$(GOOS)           Target OS (default: current)"
	@echo "  GOARCH=$(GOARCH)       Target architecture (default: current)"
