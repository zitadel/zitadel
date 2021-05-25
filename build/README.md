
# Development

## Prerequisite

- Buildkit compatible docker installation

## Generate Proto Clients

### Angular

This command generates the grpc stub for angular into the folder console/src/app/proto/generated for local development

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target npm-copy -o .
```

### Go

With this command you can generate the stub for golang into the zitadel dir 

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target go-copy -o .
```

## Run

### Run Console


#### API's from zitadel.dev

This uses zitadel.dev as API-host. If you are outside of CAOS use zitadel.ch (//TODO: how to set up a project for console) or run the entire system locally (Fullstack including database).

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml up backend-run frontend-run
```

### Run backend

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
    && docker compose -f ./build/local/docker-compose-dev.yml up -d db \
    && docker compose -f ./build/local/docker-compose-dev.yml up --exit-code-from db-migrations db-migrations \
    && docker compose -f ./build/local/docker-compose-dev.yml up --exit-code-from backend-setup backend-setup \
    && docker compose -f ./build/local/docker-compose-dev.yml up backend-run
```

### Fullstack including database

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
    && docker compose -f ./build/local/docker-compose-dev.yml up -d db \
    && docker compose -f ./build/local/docker-compose-dev.yml up --exit-code-from db-migrations db-migrations \
    && docker compose -f ./build/local/docker-compose-dev.yml up --exit-code-from backend-setup backend-setup \
    && docker compose -f ./build/local/docker-compose-dev.yml up backend-run frontend-local-run
```

## Production Build

This can also be run locally!

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --build-arg ENV=prod
```
