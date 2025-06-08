LOGIN_DEPENDENCIES_TAG ?= "zitadel-login-dependencies:local"
LOGIN_IMAGE_TAG ?= "zitadel-login:local"
CORE_MOCK_TAG ?= "zitadel-core-mock:local"
LOGIN_INTEGRATION_TESTSUITE_TAG ?= "zitadel-login-integration-testsuite:local"

XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

.PHONY: login-help
login-help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  login-help              - Show this help message."
	@echo "  login-lint              - Run linting and formatting checks. FORCE=true prevents skipping."
	@echo "  login-unit              - Run unit tests. FORCE=true prevents skipping."
	@echo "  login-integration       - Run integration tests. FORCE=true prevents skipping."
	@echo "  login-standalone-build  - Build the docker image for production login containers."
	@echo "  login-quality           - Run all quality checks (login-lint, unit, integration)."
	@echo "  login-ci                - Run all CI tasks. Run it with the -j flag to parallelize: make -j ci."
	@echo "  show-cache-keys         - Show all cache keys with image ids and exit codes."
	@echo "  clean-cache-keys        - Remove all cache keys."


login-lint-run: login-dependencies
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) lint
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) format --check

.PHONY: login-lint
login-lint:
	./scripts/run_or_skip.sh login-lint-run $(LOGIN_DEPENDENCIES_TAG)

login-unit-run: login-dependencies
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) test:unit

.PHONY: login-unit
login-unit:
	./scripts/run_or_skip.sh login-unit-run $(LOGIN_DEPENDENCIES_TAG)

login-integration-run: login-standalone-build core-mock-build login-integration-testsuite-build
	docker compose --file ./apps/login-integration-testsuite/docker-compose.yaml run --rm integration-testsuite

.PHONY: login-integration
login-integration:
	./scripts/run_or_skip.sh login-integration-run '$(LOGIN_IMAGE_TAG);$(LOGIN_INTEGRATION_TESTSUITE_TAG);$(CORE_MOCK_TAG)'

.PHONY: login-quality
login-quality: core-mock-build login-quality-after-build
login-quality-after-build: login-lint login-unit login-integration
	@:

.PHONY: login-ci
login-ci: core-mock-build login-ci-after-build
login-ci-after-build: login-quality-after-build login-standalone-build
	@:

login-dependencies:
	docker buildx bake login-dependencies --set login-dependencies.tags=$(LOGIN_DEPENDENCIES_TAG);

.PHONY: login-standalone-build
login-standalone-build:
	docker buildx bake login-standalone --set login-standalone.tags=$(LOGIN_IMAGE_TAG);

core-mock-build:
	docker buildx bake core-mock --set core-mock.tags=$(CORE_MOCK_TAG);

login-integration-testsuite-build: login-dependencies
	docker buildx bake login-integration-testsuite --set login-integration-testsuite.tags=$(LOGIN_INTEGRATION_TESTSUITE_TAG)

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
