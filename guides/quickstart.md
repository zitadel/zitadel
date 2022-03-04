# Quickstart with docker compose

We strongly recommend to go from top to bottom in this guide so all commands work seemlessly.

## Prerequisites

The commands in this guide are known to work with the following prerequisites: 

### Resources

* CPU's: 2
* Memory: 4Gb
* Free disk space: 16G

### Operating Systems
* Ubuntu, Version 18.04
* macOS Monterey, Version 12.2.1

### Dependencies
* Docker Community Editition, Version 20.10.12
* [Compose V2]((https://docs.docker.com/compose/cli-command/), Version 2.2.3

### Environment Variables

For working docker compose runs, you need to export some environment variables.

```bash
$ # In order to use Docker and Compose with buildkit enabled, export two environment variables for your current shell
$ export DOCKER_BUILDKIT=1 
$ export COMPOSE_DOCKER_CLI_BUILD=1

$ # in order to run containers as the currently logged in user, export his user and group ids
$ export UID=$(id -u) 
$ export GID=$(id -g)
```

## Starting ZITADEL

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

You can connect to [ZITADEL on localhost:4200](http://localhost:4200) after the frontend compiled successfully. Initially it takes several minutes to start all containers.

```bash
$ docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend up --detach
```

## Developing ZITADEL

If you want to make changes to ZITADEL, we recommend running the end-to-end tests against it. 

### Prerequisites

Additionally to the prerequsites described [above](#prerequisites), the end-to-end tests are known to work with the following dependencies:

* NodeJS, Version 14.17.6
* NPM, Version 6.14.15

### Developing the Backend

```bash
$ # Make changes to the backend, then rebuild and redeploy it 
$ docker compose -f ../build/local/docker-compose-local.yml up -d --no-deps --build backend-run

$ # Change to the console directory
$ cd ./console

$ # Install test suite dependencies
$ npm install --only dev

$ # Run all end-to-end tests
$ npm run e2e

$ # Or open the end-to-end test suite interactively
$ npm run e2e:open
```

### Developing the Frontend

You can switch to `ng serve` for better developer experience.

```
$ # Stop the frontend container
$ docker stop local-frontend-run-1

$ # Change to the console directory
$ cd ./console

$ # Install dependencies
$ npm install

$ # Run the local server
$ ng serve

$ # Run all end-to-end tests
$ npm run e2e

$ # Or open the end-to-end test suite interactively
$ npm run e2e:open

$ # If you want, you can stop `ng serve`, rebuild, and rerun the frontend in a new container
$ docker compose -f ./build/local/docker-compose-local.yml --profile frontend up --detach --build
```


### Regenerating gRPC Stubs

When you created your environment using docker compose, the stubs were already initially generated. If you need to change .proto files, ensure you regenerate the stubs using the following commands

```bash
$ # Backend
$ docker compose -f ./build/local/docker-compose-local.yml up -d --no-deps --build go-copy

$ # Frontend
$ docker compose -f ./build/local/docker-compose-local.yml up -d --no-deps --build npm-copy
```


## FAQ

### Initial login credentials

**username**: `zitadel-admin@caos-ag.zitadel.ch`

**password**: `Password1!`  

### Troubleshooting

#### Failing End-to-End Tests

The command `npm run e2e` prints `no such service: db` and the end-to-end test fail with many 401 responses. Make sure you have the docker compose plugin version installed as described [above](#Dependencies)

#### database-migrations don't start or stop without exit code

You can simply restart the database with the following command:

```bash
$ docker compose -f ./build/local/docker-compose-local.yml restart db
```

#### API call's block and don't return any response

The problem is that the database has a connection issues. You can simply restart the database with the following command:

```bash
$ docker compose -f ./build/local/docker-compose-local.yml restart db
```

### Destroy your local development environment

```bash
$ docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend rm
```

