
## Prerequisites

- Buildkit

## Local Dev

## Angular Generate Proto Stub

```Bash
DOCKER_BUILDKIT=1 docker build -f build/Dockerfile . -t zitadel:local --target npm-copy -o console/src/app/proto/generated
```

## Go Generate Proto Stub

```Bash
DOCKER_BUILDKIT=1 docker build -f build/Dockerfile . -t zitadel:local --target go-copy -o pkg
```

### Angular Run

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-dev.yml up --build angular
```

### Go Run

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-dev.yml up --build  go
```

### Go Debug

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose-debug.yml up --build  go
```

### Go, Dd

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f build/docker-compose.yml up go db --build
```

## Production Build

This can also be run locally!

```Bash
DOCKER_BUILDKIT=1 docker build -f build/Dockerfile . -t zitadel:local --build-arg ENV=prod
```
