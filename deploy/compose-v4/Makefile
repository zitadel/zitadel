ENV_FILE ?= .env
COMPOSE ?= docker compose --env-file $(ENV_FILE)
BASE_FILE ?= -f docker-compose.yml
OVERRIDES ?=
PROFILES ?=

.PHONY: up down logs reset config release-bump

up:
	$(COMPOSE) $(BASE_FILE) $(OVERRIDES) $(PROFILES) up -d --wait

down:
	$(COMPOSE) $(BASE_FILE) $(OVERRIDES) $(PROFILES) down

logs:
	$(COMPOSE) $(BASE_FILE) $(OVERRIDES) $(PROFILES) logs -f --tail=200

reset:
	$(COMPOSE) $(BASE_FILE) $(OVERRIDES) $(PROFILES) down -v --remove-orphans

config:
	$(COMPOSE) $(BASE_FILE) $(OVERRIDES) $(PROFILES) config

release-bump:
	@echo "TODO: wire Nx release automation after nx release PR merges."
	@echo "Expected future command shape: nx release version --specifier <patch|minor|major>"
