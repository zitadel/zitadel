export LOGIN_IMAGE_TAG ?= zitadel-login:local
LOGIN_LINT_TAG ?= zitadel-login-lint:local
LOGIN_DEPENDENCIES_TAG ?= zitadel-login-dependencies:local
LOGIN_TEST_UNIT_TAG ?= zitadel-login-lint:local
export CORE_MOCK_TAG ?= zitadel-core-mock:local
export LOGIN_TEST_INTEGRATION_TAG ?= zitadel-login-test-integration:local
export LOGIN_TEST_ACCEPTANCE_SETUP_TAG := zitadel-login-test-acceptance-setup:local
export LOGIN_TEST_ACCEPTANCE_POSTGRES_TAG := postgres:17.0-alpine3.19
export LOGIN_TEST_ACCEPTANCE_GOLANG_TAG := golang:1.24-alpine
export ZITADEL_IMAGE_TAG ?= ghcr.io/zitadel/zitadel:latest

XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

.PHONY: login-help
login-help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  login-help              - Show this help message."
	@echo "  login-lint              - Run linting and formatting checks. FORCE=true prevents skipping."
	@echo "  login-test-unit         - Run unit tests. FORCE=true prevents skipping."
	@echo "  login-test-integration  - Run integration tests. FORCE=true prevents skipping."
	@echo "  login-standalone-build  - Build the docker image for production login containers."
	@echo "  login-quality           - Run all quality checks (login-lint, login-unit, login-integration)."
	@echo "  login-ci                - Run all CI tasks. Run it with the -j flag to parallelize: make -j ci."
	@echo "  show-cache-keys         - Show all cache keys with image ids and exit codes."
	@echo "  clean-cache-keys        - Remove all cache keys."


login-lint-run:
	docker run --rm $(LOGIN_LINT_TAG) lint
	docker run --rm $(LOGIN_LINT_TAG) format --check

.PHONY: login-lint
login-lint: login-lint-build
	./scripts/run_or_skip.sh login-lint-run $(LOGIN_LINT_TAG)

login-test-unit-run:
	docker run --rm $(LOGIN_TEST_UNIT_TAG) test:unit:standalone

.PHONY: login-test-unit
login-test-unit: login-test-unit-build
	./scripts/run_or_skip.sh login-test-unit-run $(LOGIN_TEST_UNIT_TAG)

login-test-integration-run:
	docker compose --file ./apps/login-test-integration/docker-compose.yaml run --rm login-test-integration

.PHONY: login-test-integration
login-test-integration: login-standalone-build login-test-integration-build
	./scripts/run_or_skip.sh login-test-integration-run "$(LOGIN_IMAGE_TAG);$(CORE_MOCK_TAG);$(LOGIN_TEST_INTEGRATION_TAG)"

login-test-acceptance-run:
	docker compose --file ./apps/login-test-acceptance/saml/docker-compose.yaml up --detach samlsp
	docker compose --file ./apps/login-test-acceptance/oidc/docker-compose.yaml up --detach oidcrp
	docker compose --file ./apps/login-test-acceptance/docker-compose.yaml run login-test-acceptance

login-test-acceptance: login-standalone-build login-test-acceptance-build
	./scripts/run_or_skip.sh login-test-acceptance-run "$(LOGIN_IMAGE_TAG);$(LOGIN_TEST_ACCEPTANCE_SETUP_TAG);$(LOGIN_TEST_ACCEPTANCE_POSTGRES_TAG);$(LOGIN_TEST_ACCEPTANCE_GOLANG_TAG)"

.PHONY: login-quality
login-quality: login-lint login-test-unit login-test-integration
	@:

.PHONY: login-ci
login-ci: login-quality login-standalone-build
	@:

login-dependencies-build:
	docker buildx bake login-dependencies --set login-dependencies.tags=$(LOGIN_DEPENDENCIES_TAG);

login-lint-build:
	docker buildx bake login-lint --set login-lint.tags=$(LOGIN_LINT_TAG);

login-test-unit-build:
	docker buildx bake login-test-unit --set login-test-unit.tags=$(LOGIN_TEST_UNIT_TAG);

login-test-integration-build:
	docker buildx bake core-mock --set core-mock.tags=$(CORE_MOCK_TAG);
	docker buildx bake login-test-integration --set login-test-integration.tags=$(LOGIN_TEST_INTEGRATION_TAG)

login-test-acceptance-build:
	# TODO: Prebuild sink, saml and oidc
	docker buildx bake --pull --file apps/login-test-acceptance/docker-compose.yaml --set setup.context=apps/login-test-acceptance

.PHONY: login-standalone-build
login-standalone-build:
	docker buildx bake login-standalone --set login-standalone.tags=$(LOGIN_IMAGE_TAG);

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
