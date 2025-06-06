LOGIN_BASE_TAG ?= "zitadel-login-base:local"
CORE_MOCK_TAG ?= "zitadel-core-mock:local"
XDG_CACHE_HOME ?= $(HOME)/.cache
CACHE_DIR ?= $(XDG_CACHE_HOME)/zitadel-make

.PHONY: help
help:
	@echo "Makefile for the login service"
	@echo "Available targets:"
	@echo "  help              - Show this help message"
	@echo "  lint              - Run linting and formatting checks"
	@echo "  lint-force        - Force run linting and formatting checks"
	@echo "  unit              - Run unit tests"
	@echo "  unit-force        - Force run unit tests"
	@echo "  integration       - Run integration tests"
	@echo "  integration-force - Force run integration tests"
	@echo "  login-image       - Build the login image"
	@echo "  quality           - Run all quality checks (lint, unit, integration)"
	@echo "  ci                - Run all CI tasks. Run it with the -j flag to parallelize. make -j ci"

.PHONY: lint-force
lint-force:
	docker run --rm $(LOGIN_BASE_TAG) lint
	docker run --rm $(LOGIN_BASE_TAG) format --check

.PHONY: lint
lint:
	$(call run_or_skip,lint-force,lint,$(LOGIN_BASE_TAG))

unit-run: login-base
	docker run --rm $(LOGIN_BASE_TAG) test:unit

.PHONY: unit-force
unit-force:
	docker run --rm $(LOGIN_BASE_TAG) test:unit

.PHONY: unit
unit:
	$(call run_or_skip,unit-force,unit,$(LOGIN_BASE_TAG))

.PHONY: integration-force
integration-force:
	docker run --rm $(CORE_MOCK_TAG) test:integration

.PHONY: integration
integration:
	$(call run_or_skip,integration-force,integration,$(CORE_MOCK_TAG))

.PHONY: login-image
login-image:
	docker buildx bake login-image

.PHONY: quality
quality: lint unit integration

.PHONY: ci
ci: core-mock ci-after-build
ci-after-build: quality login-image
	@:

login-base:
	docker buildx bake login-base --set login-base.tags=$(LOGIN_BASE_TAG);

core-mock:
	docker buildx bake core-mock --set login-base.tags=$(CORE_MOCK_TAG);

.PHONY: clean-cache
clean-cache:
	@echo "Removing cache directory: $(CACHE_DIR)"
	@rm -rf "$(CACHE_DIR)"

.PHONY: show-cache
show-cache:
	@echo "Showing cached digests and exit codes in $(CACHE_DIR):"
	@find "$(CACHE_DIR)" -type f 2>/dev/null | while read file; do \
		echo "$$file: $$(cat $$file)"; \
	done

# run_or_skip: runs a task only if the Docker image has changed and caches the result
# $(1): Taskname (e.g. "lint-force")
# $(2): Cache-ID (e.g. "lint")
# $(3): Docker-Image (e.g. "zitadel-login-base:local")
define run_or_skip
	@digest_file="$(CACHE_DIR)/$(2).$(3)"; \
	mkdir -p $(CACHE_DIR); \
	if [ -f "$$digest_file" ]; then \
		digest_before=$$(cut -d',' -f1 "$$digest_file"); \
		status_before=$$(cut -d',' -f2 "$$digest_file"); \
	else \
		digest_before=""; \
		status_before=1; \
	fi; \
	current_digest=$$(docker image inspect $(3) --format='{{.Id}}'); \
	if [ "$$digest_before" = "$$current_digest" ]; then \
		echo "Skipping $(1) – image unchanged, returning cached status $$status_before"; \
		exit $$status_before; \
	else \
		echo "Running $(1)..."; \
		$(MAKE) $(1); \
		status=$$?; \
		echo "$$current_digest,$$status" > "$$digest_file"; \
		exit $$status; \
	fi
endef
