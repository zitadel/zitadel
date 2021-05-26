# Startup ZITADEL

## Prerequesits

* docker
* go (for backend local development)
* ng (for frontend local development)


## On system

### Keys file

Generates the required keys for cryptography.

```bash
docker build --target copy_keys -f build/Dockerfile.dev . -o .keys
```

### env variables

Default env variables are provided in [this .env-file](build/local/local.env)

## Pre steps

### generate code

```bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target go-copy -o .
```

### start database (cockroach)

The following command creates the dbms and sets up the database structure

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml up db db-migrations
```

## setup ZITADEL

You can use your local debugger or you can setup with a docker compose

Make sure that the database is ready and migrations are finished

### local

`go run cmd/zitadel/main.go -setup-files=cmd/zitadel/setup.yaml -setup-files=cmd/zitadel/system-defaults.yaml -setup-files=cmd/zitadel/authz.yaml setup`

### docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml up --exit-code-from db-migrations db-migrations
```

## start backend

You can use your local debugger or you can start the backend with a docker compose

Make sure that the database is ready, migrations are finished and setup ended successfully

### local

`go run -console=false -localDevMode=true -config-files=cmd/zitadel/startup.yaml -config-files=cmd/zitadel/system-defaults.yaml -config-files=cmd/zitadel/authz.yaml start`

### docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml up backend-run
```

## start frontend

Make sure that the defined backend is serving traffic.  
Both options support hot reloading.


### local

```bash
cd console
ng serve --host localhost
```

### docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml up frontend-local-run
```




zitadel-admin@caos-ag.zitadel.ch
Password1!