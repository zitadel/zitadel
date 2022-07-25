#!/bin/bash

COMPOSE_DOCKER_CLI_BUILD=1 docker compose --file ${projectRoot}/e2e/docker-compose-workdir.yaml --file ${projectRoot}/docs/docs/guides/installation/run/docker-compose.yaml --file ${projectRoot}/e2e/docker-compose-overwrite.yaml --env-file ${projectRoot}/e2e/compose.env "$@"
