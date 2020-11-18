
# Development

## Prerequisite

- Buildkit compatible docker installation

## Generate Proto Clients

### Angular

This command generates the grpc stub for angular into the folder console/src/app/proto/generated for local development

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target npm-copy -o console/src/app/proto/generated
```

### Go

With this command you can generate the stub for golang into the correct dir pkg/

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target go-copy -o pkg
```

## Run

### Run Angular

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-dev.yml up --build angular
```

### Run Go

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-dev.yml up --build  go
```

### Fullstack including database

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-dev.yml up --build
```

## Debug

### Debug Go

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-debug.yml up --build  go
```

## Production Build

This can also be run locally!

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --build-arg ENV=prod
```
