# Development

You should stay in the ZITADEL root directory to execute the statements in the following chapters.

## Prerequisite

- Buildkit compatible docker installation

### env variables

Default env variables are provided in [this .env-file](../build/local/local.env)

## Generate required files

This part is relevant if you start the backend or console without docker compose.

### Console

This command generates the grpc stub for console into the folder console/src/app/proto/generated for local development.

```bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:gen-fe --target npm-copy -o .
```

### Backend

With this command you can generate the stub for the backend.

```bash
# generates grpc stub
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:gen-be --target go-copy -o .
# generates keys for cryptography
DOCKER_BUILDKIT=1 docker build --target copy_keys -f build/Dockerfile.dev . -o .keys
```

## Run

### Initialise data

Used if you want to setup the database and load the initial data.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend -p zitadel up
```

You can stop as soon as db-migrations AND backend-setup returned with exit code 0.

### Initialise frontend

Used to set the client id of the console This step is for local development. If you don't work with a local backend you have to set the client id manually.

You must [initialise the data](###-Initialise-data)) first.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile backend --profile init-frontend -p zitadel up --exit-code-from client-id
```

The command exists as soon as the client id is set.

### Run database

Used if you want to run the backend/console locally and only need the database. It's recommended to [initialise the data](###-Initialise-data) first.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-local.yml --profile database -p zitadel up
```

**On apple silicon:**
Restart the command (second terminal `docker restart zitadel-<SERVICE_NAME>_1`) if `db` logs `qemu: uncaught target signal 11 (Segmentation fault) - core dumped` or no logs are written from `db-migrations`.

### Run Console

The console service is configured for hot reloading. You can also use docker compose for local development.

If you don't use the backend from local you have to configure [the environment.json](../build/local/environment.json) manually.

If you use the local backend ensure that you run that you have [set the correct client id](###-Initialise-frontend).

#### Docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-local.yml --profile frontend -p zitadel up
```

### Run backend

Used if you want to run the backend locally. It's recommended to [initialise the data](###-Initialise-data) first.

#### Docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml --profile database --profile backend -p zitadel up
```

#### Local

##### Export environment variables

```bash
# exports all default env variables
while read line; do
    if [[ $line != #* ]] && [[ ! -z $line ]]; then
        export $line
    fi
done < build/local/local.env
```

##### Start ZITADEL

```bash
# starts zitadel with default config files
go run cmd/zitadel/main.go -console=false -localDevMode=true -config-files=cmd/zitadel/startup.yaml -config-files=cmd/zitadel/system-defaults.yaml -config-files=cmd/zitadel/authz.yaml start
```

If you want to run your backend locally and the frontend by docker compose you have to replace the following variables:

[docker compose yaml](../build/local/docker-compose-local.yml):

```yaml
service:
  client-id:
    environment:
      - HOST=backend-run
  grpc-web-gateway:
    environment:
      - BKD_HOST=backend-run
```

with

```yaml
service:
  client-id:
    environment:
      - HOST=host.docker.internal
  grpc-web-gateway:
    environment:
      - BKD_HOST=host.docker.internal
```

##### Setup ZITADEL

```bash
# starts zitadel with default config files
go run cmd/zitadel/main.go -setup-files=cmd/zitadel/setup.yaml -setup-files=cmd/zitadel/system-defaults.yaml -setup-files=cmd/zitadel/authz.yaml setup
```

## Initial login credentials

**username**: `zitadel-admin@caos-ag.zitadel.ch`

**password**: `Password1!`