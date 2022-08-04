#!/bin/bash

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

COMPOSE_DOCKER_CLI_BUILD=1 docker compose --file ${SCRIPTPATH}/docker-compose-workdir.yaml --file ${SCRIPTPATH}/../docs/docs/guides/deploy/docker-compose.yaml --file ${SCRIPTPATH}/docker-compose-overwrite.yaml "$@"
