XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

LOGIN_DIR ?= ./
LOGIN_BAKE_CLI ?= docker buildx bake
LOGIN_BAKE_CLI_WITH_ARGS := $(LOGIN_BAKE_CLI) --file $(LOGIN_DIR)docker-bake.hcl --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml
LOGIN_BAKE_CLI_ADDITIONAL_ARGS ?=
LOGIN_BAKE_CLI_WITH_ARGS += $(LOGIN_BAKE_CLI_ADDITIONAL_ARGS)

export COMPOSE_BAKE=true
export UID := $(id -u)
export GID := $(id -g)

export LOGIN_TEST_ACCEPTANCE_BUILD_CONTEXT := $(LOGIN_DIR)apps/login-test-acceptance

export DOCKER_METADATA_OUTPUT_VERSION ?= local
export LOGIN_TAG ?= login:${DOCKER_METADATA_OUTPUT_VERSION}
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
export ZITADEL_TAG ?= ghcr.io/zitadel/zitadel:latest
export LOGIN_CORE_MOCK_TAG := login-core-mock:${DOCKER_METADATA_OUTPUT_VERSION}

login_help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  login_help              - Show this help message."
	@echo "  login_quality           - Run all quality checks (login_lint, login_test_unit, login_test_integration, login_test_acceptance)."
	@echo "  login_standalone_build  - Build the docker image for production login containers."
	@echo "  login_lint              - Run linting and formatting checks. IGNORE_RUN_CACHE=true prevents skipping."
	@echo "  login_test_unit         - Run unit tests. Tests without any dependencies. IGNORE_RUN_CACHE=true prevents skipping."
	@echo "  login-test_integration  - Run integration tests. Tests a login production build against a mocked Zitadel core API. IGNORE_RUN_CACHE=true prevents skipping."
	@echo "  login_test_acceptance   - Run acceptance tests. Tests a login production build with a local Zitadel instance behind a reverse proxy. IGNORE_RUN_CACHE=true prevents skipping."
	@echo "  typescript_generate	 - Generate TypeScript client code from Protobuf definitions."
	@echo "  show_run_caches         - Show all run caches with image ids and exit codes."
	@echo "  clean_run_caches        - Remove all run caches."


login_lint:
	@echo "Running login linting and formatting checks"
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-lint

login_test_unit:
	@echo "Running login unit tests"
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-test-unit

login_test_integration_build:
	@echo "Building login integration test environment with the local core mock image"
	$(LOGIN_BAKE_CLI_WITH_ARGS) core-mock login-test-integration login-standalone --load

login_test_integration_dev: login_test_integration_cleanup
	@echo "Starting login integration test environment with the local core mock image"
	$(LOGIN_BAKE_CLI_WITH_ARGS) core-mock && docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml run --service-ports --rm core-mock

login_test_integration_run: login_test_integration_cleanup
	@echo "Running login integration tests"
	docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml run --rm integration

login_test_integration_cleanup:
	@echo "Cleaning up login integration test environment"
	docker compose --file $(LOGIN_DIR)apps/login-test-integration/docker-compose.yaml down --volumes

login_test_integration: login_test_integration_build
	$(LOGIN_DIR)scripts/run_or_skip.sh login_test_integration_run \
	"$(LOGIN_TAG) \
	$(LOGIN_CORE_MOCK_TAG) \
	$(LOGIN_TEST_INTEGRATION_TAG)"

login_test_acceptance_build_bake:
	@echo "Building login test acceptance images as defined in the docker-bake.hcl"
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-test-acceptance login-standalone --load

login_test_acceptance_build_compose:
	@echo "Building login test acceptance images as defined in the docker-compose.yaml"
	$(LOGIN_BAKE_CLI_WITH_ARGS) --load setup sink

# login_test_acceptance_build is overwritten by the login_dev target in zitadel/zitadel/Makefile
login_test_acceptance_build: login_test_acceptance_build_compose login_test_acceptance_build_bake

login_test_acceptance_run: login_test_acceptance_cleanup
	@echo "Running login test acceptance tests"
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose-ci.yaml run --rm --service-ports acceptance

login_test_acceptance_cleanup:
	@echo "Cleaning up login test acceptance environment"
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose-ci.yaml down --volumes

login_test_acceptance: login_test_acceptance_build
	$(LOGIN_DIR)scripts/run_or_skip.sh login_test_acceptance_run \
		"$(LOGIN_TAG) \
  		$(ZITADEL_TAG) \
  		$(POSTGRES_TAG) \
  		$(GOLANG_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SETUP_TAG) \
  		$(LOGIN_TEST_ACCEPTANCE_SINK_TAG)"

login_test_acceptance_setup_env: login_test_acceptance_build_compose login_test_acceptance_cleanup
	@echo "Setting up the login test acceptance environment and writing the env.test.local file"
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml run setup

login_test_acceptance_setup_dev:
	@echo "Starting the login test acceptance environment with the local zitadel image"
	docker compose --file $(LOGIN_DIR)apps/login-test-acceptance/docker-compose.yaml up --no-recreate zitadel traefik sink

login_quality: login_lint login_test_unit login_test_integration
	@echo "Running login quality checks: lint, unit tests, integration tests"

login_standalone_build:
	@echo "Building the login standalone docker image with tag: $(LOGIN_TAG)"
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-standalone --load

login_standalone_out:
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-standalone-out

typescript_generate:
	@echo "Generating TypeScript client and writing to local $(LOGIN_DIR)packages/zitadel-proto"
	$(LOGIN_BAKE_CLI_WITH_ARGS) login-typescript-proto-client-out

clean_run_caches:
	@echo "Removing cache directory: $(CACHE_DIR)"
	rm -rf "$(CACHE_DIR)"

show_run_caches:
	@echo "Showing run caches with docker image ids and exit codes in $(CACHE_DIR):"
	@find "$(CACHE_DIR)" -type f 2>/dev/null | while read file; do \
		echo "$$file: $$(cat $$file)"; \
	done

