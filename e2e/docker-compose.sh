#!/bin/bash

COMPOSE_DOCKER_CLI_BUILD=1 docker compose --file ${projectRoot}/docs/docs/guides/installation/run/docker-compose.yaml --file ${projectRoot}/e2e/docker-compose-overwrite.yaml "$@"
