LOGIN_DEPENDENCIES_TAG ?= "zitadel-login-dependencies:local"
LOGIN_IMAGE_TAG ?= "zitadel-login:local"
CORE_MOCK_TAG ?= "zitadel-core-mock:local"
LOGIN_INTEGRATION_TESTSUITE_TAG ?= "zitadel-login-integration-testsuite:local"
CORE_MOCK_CONTAINER_NAME ?= zitadel-mock-grpc-server
LOGIN_CONTAINER_NAME ?= zitadel-login

XDG_CACHE_HOME ?= $(HOME)/.cache
export CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

.PHONY: help
help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  help              	     - Show this help message"
	@echo "  login                   - Start the login service"
	@echo "  login-lint              - Run linting and formatting checks"
	@echo "  login-lint-force        - Force run linting and formatting checks"
	@echo "  login-unit              - Run unit tests"
	@echo "  login-unit-force        - Force run unit tests"
	@echo "  login-integration       - Run integration tests"
	@echo "  login-integration-force - Force run integration tests"
	@echo "  login-standalone        - Build the docker image for production login containers"
	@echo "  login-quality           - Run all quality checks (login-lint, unit, integration)"
	@echo "  login-ci                - Run all CI tasks. Run it with the -j flag to parallelize. make -j ci"
	@echo "  show-cache-keys         - Show all cache keys with image ids and exit codes"
	@echo "  clean-cache-keys        - Remove all cache keys"
	@echo "  core-mock               - Start the core mock server"
	@echo "  core-mock-stop          - Stop the core mock server"


.PHONY: login-lint-force
login-lint-force: login-dependencies
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) lint
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) format --check

.PHONY: login-lint
login-lint:
	./scripts/run_or_skip.sh login-lint-force $(LOGIN_DEPENDENCIES_TAG)

.PHONY: login-unit-force
login-unit-force: login-dependencies
	docker run --rm $(LOGIN_DEPENDENCIES_TAG) test:unit

.PHONY: login-unit
login-unit:
	./scripts/run_or_skip.sh login-unit-force $(LOGIN_DEPENDENCIES_TAG)

.PHONY: login-integration-force
login-integration-force: login core-mock login-integration-testsuite
	docker run --rm $(LOGIN_INTEGRATION_TESTSUITE_TAG)
	$(MAKE) core-mock-stop

.PHONY: login-integration
login-integration:
	./scripts/run_or_skip.sh login-integration-force '$(LOGIN_DEPENDENCIES_TAG);$(CORE_MOCK_TAG);$(LOGIN_INTEGRATION_TESTSUITE_TAG)'

.PHONY: login-quality
login-quality: core-mock-build login-quality-after-build
login-quality-after-build: login-lint login-unit login-integration
	@:

.PHONY: login-ci
login-ci: core-mock-build login-ci-after-build
login-ci-after-build: login-quality-after-build login-standalone
	@:

login-dependencies:
	docker buildx bake login-dependencies --set login-dependencies.tags=$(LOGIN_DEPENDENCIES_TAG);

.PHONY: login-standalone
login-standalone:
	docker buildx bake login-standalone --set login-standalone.tags=$(LOGIN_IMAGE_TAG);

.PHONY: login
login: login-standalone login-stop
	docker run --detach --rm --name $(LOGIN_CONTAINER_NAME) --publish 3000:3000 $(LOGIN_IMAGE_TAG)

login-stop:
	docker rm --force $(LOGIN_CONTAINER_NAME) 2>/dev/null || true

core-mock-build:
	docker buildx bake core-mock --set core-mock.tags=$(CORE_MOCK_TAG);

login-integration-testsuite: login-dependencies
	docker buildx bake login-integration-testsuite --set login-integration-testsuite.tags=$(LOGIN_INTEGRATION_TESTSUITE_TAG)

.PHONY: core-mock
core-mock: core-mock-build core-mock-stop
	docker run --detach --rm --name $(CORE_MOCK_CONTAINER_NAME) --publish 22221:22221 --publish 22222:22222 $(CORE_MOCK_TAG)

.PHONY: core-mock-stop
core-mock-stop:
	docker rm --force $(CORE_MOCK_CONTAINER_NAME) 2>/dev/null || true

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
