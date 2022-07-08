#!/bin/bash

COMPOSE_DOCKER_CLI_BUILD=1 docker compose --file ./docs/docs/guides/installation/run/docker-compose.yaml --file ./e2e/docker-compose-overwrite.yaml "$@"
