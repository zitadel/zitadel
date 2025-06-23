XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

LOGIN_DIR ?= ./
LOGIN_BAKE_CLI ?= docker buildx bake
LOGIN_BAKE_CLI_WITH_COMMON_ARGS := $(LOGIN_BAKE_CLI) --file $(LOGIN_DIR)docker-bake.hcl --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml
LOGIN_BAKE_CLI_ADDITIONAL_ARGS ?=
LOGIN_BAKE_CLI_WITH_COMMON_ARGS += $(LOGIN_BAKE_CLI_ADDITIONAL_ARGS)

export COMPOSE_BAKE=true
export UID := $(id -u)
export GID := $(id -g)

export LOGIN_TEST_ACCEPTANCE_BUILD_CONTEXT := $(LOGIN_DIR)apps/login-test-acceptance

export DOCKER_METADATA_OUTPUT_VERSION ?= local
export LOGIN_TAG := login:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_UNIT_TAG := login-test-unit:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_INTEGRATION_TAG := login-test-integration:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_TAG := login-test-acceptance:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_SETUP_TAG := login-test-acceptance-setup:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_SINK_TAG := login-test-acceptance-sink:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_OIDCRP_TAG := login-test-acceptance-oidcrp:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_OIDCOP_TAG := login-test-acceptance-oidcop:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_SAMLSP_TAG := login-test-acceptance-samlsp:${DOCKER_METADATA_OUTPUT_VERSION}
export LOGIN_TEST_ACCEPTANCE_SAMLIDP_TAG := login-test-acceptance-samlidp:${DOCKER_METADATA_OUTPUT_VERSION}
export POSTGRES_TAG := postgres:17.0-alpine3.19
export GOLANG_TAG := golang:1.24-alpine
export ZITADEL_TAG ?= ghcr.io/zitadel/zitadel:v3.3.0
export CORE_MOCK_TAG := core-mock:${DOCKER_METADATA_OUTPUT_VERSION}

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
	@echo "  show-run-caches         - Show all run caches with image ids and exit codes."
	@echo "  clean-run-caches        - Remove all run caches."

login-lint:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) login-lint

login-test-unit:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) login-test-unit

login-test-integration-build:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) core-mock login-test-integration login-standalone

login-test-integration-dev: login-test-integration-cleanup
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) core-mock && docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml run --service-ports --rm core-mock

login-test-integration-run: login-test-integration-cleanup
	docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml run --rm integration

login-test-integration-cleanup:
	docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml down --volumes

.PHONY: login-test-integration
login-test-integration: login-test-integration-build
	$(LOGIN_DIR)scripts/run_or_skip.sh login-test-integration-run \
	"$(LOGIN_TAG) \
	$(CORE_MOCK_TAG) \
	$(LOGIN_TEST_INTEGRATION_TAG)"

login-test-acceptance-build-bake:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) login-test-acceptance login-standalone

login-test-acceptance-build-compose:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) --load setup sink

login-test-acceptance-build: login-test-acceptance-build-compose login-test-acceptance-build-bake
	@:

login-test-acceptance-dev: login-test-acceptance-build-compose login-test-acceptance-cleanup
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml up zitadel setup traefik setup sink

login-test-acceptance-run: login-test-acceptance-cleanup
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose-ci.yaml run --rm --service-ports acceptance

login-test-acceptance-cleanup:
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose-ci.yaml down --volumes

login-test-acceptance: login-test-acceptance-build
	$(LOGIN_DIR)scripts/run_or_skip.sh login-test-acceptance-run \
		"$(LOGIN_TAG) \
  		$(ZITADEL_TAG) \
  		$(POSTGRES_TAG) \
  		$(GOLANG_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SETUP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SINK_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_OIDCRP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SAMLSP_TAG)"

.PHONY: login-quality
login-quality: login-lint login-test-unit login-test-integration
	@:

.PHONY: login-standalone-build
login-standalone-build:
	$(LOGIN_BAKE_CLI_WITH_COMMON_ARGS) --load login-standalone

login-standalone-build-tag:
	@echo -n "$(LOGIN_TAG)"

.PHONY: clean-run-caches
clean-run-caches:
	@echo "Removing cache directory: $(CACHE_DIR)"
	rm -rf "$(CACHE_DIR)"

.PHONY: show-run-caches
show-run-caches:
	@echo "Showing run caches with docker image ids and exit codes in $(CACHE_DIR):"
	@find "$(CACHE_DIR)" -type f 2>/dev/null | while read file; do \
		echo "$$file: $$(cat $$file)"; \
	done
