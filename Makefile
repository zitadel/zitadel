XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

export LOGIN_TAG ?= login:local
export LOGIN_TEST_UNIT_TAG := login-test-unit:local
export LOGIN_TEST_INTEGRATION_TAG ?= login-test-integration:local
export LOGIN_TEST_ACCEPTANCE_TAG := login-test-acceptance:local
export LOGIN_TEST_ACCEPTANCE_SETUP_TAG := login-test-acceptance-setup:local
export LOGIN_TEST_ACCEPTANCE_SINK_TAG := login-test-acceptance-sink:local
export LOGIN_TEST_ACCEPTANCE_OIDCRP_TAG := login-test-acceptance-oidcrp:local
export LOGIN_TEST_ACCEPTANCE_OIDCOP_TAG := login-test-acceptance-oidcop:local
export LOGIN_TEST_ACCEPTANCE_SAMLSP_TAG := login-test-acceptance-samlsp:local
export LOGIN_TEST_ACCEPTANCE_SAMLIDP_TAG := login-test-acceptance-samlidp:local
export POSTGRES_TAG := postgres:17.0-alpine3.19
export GOLANG_TAG := golang:1.24-alpine
# TODO: use ghcr.io/zitadel/zitadel:latest
export ZITADEL_TAG ?= ghcr.io/zitadel/zitadel:02617cf17fdde849378c1a6b5254bbfb2745b164
export CORE_MOCK_TAG := core-mock:local

.PHONY: login-help
login-help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  login-help              - Show this help message."
	@echo "  login-quality           - Run all quality checks (login-lint, login-test-unit, login-test-integration, login-test-acceptance)."
	@echo "  login-standalone-build  - Build the docker image for production login containers."
	@echo "  login-lint              - Run linting and formatting checks. FORCE=true prevents skipping."
	@echo "  login-test-unit         - Run unit tests. Tests without any dependencies. FORCE=true prevents skipping."
	@echo "  login-test-integration  - Run integration tests. Tests a login production build against a mocked Zitadel core API. FORCE=true prevents skipping."
	@echo "  login-test-acceptance   - Run acceptance tests. Tests a login production build with a local Zitadel instance behind a reverse proxy. FORCE=true prevents skipping."
	@echo "  show-cache-keys         - Show all cache keys with image ids and exit codes."
	@echo "  clean-cache-keys        - Remove all cache keys."

login-lint-build:
	docker buildx bake login-lint

login-lint-run:
	docker run --rm $(LOGIN_LINT_TAG) lint
	docker run --rm $(LOGIN_LINT_TAG) format --check

.PHONY: login-lint
login-lint: login-lint-build
#	./scripts/run_or_skip.sh login-lint-run $(LOGIN_LINT_TAG)

login-test-unit-build:
	docker buildx bake login-test-unit

login-test-unit-run:
	docker run --rm $(LOGIN_TEST_UNIT_TAG) test:unit:standalone

.PHONY: login-test-unit
login-test-unit: login-test-unit-build
	./scripts/run_or_skip.sh login-test-unit-run $(LOGIN_TEST_UNIT_TAG)

login-test-integration-build:
	docker buildx bake core-mock
	docker buildx bake login-test-integration

login-test-integration-run: login-test-integration-cleanup
	docker compose --file ./apps/login-test-integration/docker-compose.yaml run --rm integration

login-test-integration-cleanup:
	docker compose --file ./apps/login-test-integration/docker-compose.yaml down --volumes

.PHONY: login-test-integration
login-test-integration: login-standalone-build login-test-integration-build
	./scripts/run_or_skip.sh login-test-integration-run \
	"$(LOGIN_TAG) \
	$(CORE_MOCK_TAG) \
	$(LOGIN_TEST_INTEGRATION_TAG)"

login-test-acceptance-build:
	COMPOSE_BAKE=true docker compose --file ./apps/login-test-acceptance/docker-compose.yaml build
	docker buildx bake login-standalone
	docker buildx bake login-test-acceptance

login-test-acceptance-run: login-acceptance-cleanup
	docker compose --file ./apps/login-test-acceptance/docker-compose.yaml run --rm --service-ports acceptance

login-acceptance-cleanup:
	docker compose --file ./apps/login-test-acceptance/docker-compose.yaml down --volumes

login-test-acceptance: login-standalone-build login-test-acceptance-build
	./scripts/run_or_skip.sh login-test-acceptance-run \
		"$(LOGIN_TAG) \
  		$(ZITADEL_TAG) \
  		$(POSTGRES_TAG) \
  		$(GOLANG_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SETUP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SINK_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_OIDCRP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_OIDCOP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SAMLSP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SAMLIDP_TAG)"

.PHONY: login-quality
login-quality: login-lint login-test-unit login-test-integration login-test-acceptance
	@:

.PHONY: login-standalone-build
login-standalone-build:
	docker buildx bake login-standalone

.PHONY: clean-cache-keys
clean-cache-keys:
	@echo "Removing cache directory: $(CACHE_DIR)"
	rm -rf "$(CACHE_DIR)"

.PHONY: show-cache-keys
show-cache-keys:
	@echo "Showing cache keys with docker image ids and exit codes in $(CACHE_DIR):"
	@find "$(CACHE_DIR)" -type f 2>/dev/null | while read file; do \
		echo "$$file: $$(cat $$file)"; \
	done
