# Start with docker compose

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs. Until the essential services are started and executed some services panic, this is expected just give it some minutes to setup the database and execute migrations.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml --profile init --profile backend --profile frontend -p zitadel up
```

# Development

You should stay in the ZITADEL root directory to execute the statements in the following chapters.

## Prerequisite

- Buildkit compatible docker installation

## Generate Proto Clients

This part is relevant if you start the backend or console without docker compose.

### Console

This command generates the grpc stub for console into the folder console/src/app/proto/generated for local development.

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:gen-fe --target npm-copy -o .
```

### Backend

With this command you can generate the stub for the backend.

```Bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:gen-be --target go-copy -o .
```

## Run

### Initialise data

Used if you want to setup the database and load the initial data.

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-dev.yml --profile database --profile init-backend -p zitadel up
```

You can stop as soon as db-migrations AND backend-setup returned with exit code 0.

### Initialise frontend

Used to set the client id of the console This step is for local development. If you don't work with a local backend you have to set the client id manually. This

You must [initialise the data](###-Initialise-data)) first.

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-dev.yml --profile database --profile backend --profile init-frontend -p zitadel up
```

You can stop as soon as client-id returned with exit code 0.

### Run database

Used if you want to run the backend/console locally and only need the database. It's recommended to [initialise the data](###-Initialise-data) first.

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-dev.yml --profile database -p zitadel up
```

### Run Console

The console service is configured for hot reloading. You can also use docker compose for local development.

If you don't use the backend from local you have to configure [the environment.json](local/environment.json) manually.

If you use the local backend ensure that you run that you have [set the correct client id](###-Initialise-frontend).

#### Docker compose

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f ./build/local/docker-compose-dev.yml --profile frontend -p zitadel up
```

### Run backend

Used if you want to run the backend locally. It's recommended to [initialise the data](###-Initialise-data) first.

#### Docker compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-dev.yml --profile database --profile backend -p zitadel up
```

#### local

```bash
# exports all default env variables
while read line; do
    if [[ $line != #* ]] && [[ ! -z $line ]]; then
        export $line
    fi
done < build/local/local.env

# starts zitadel with default config files
go run cmd/zitadel/main.go -console=false -localDevMode=true -config-files=cmd/zitadel/startup.yaml -config-files=cmd/zitadel/system-defaults.yaml -config-files=cmd/zitadel/authz.yaml start
```

# Production Build

This can also be run locally!

```bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --build-arg ENV=prod
```