# Quickstart with docker compose

We strongly recommend to go from top to bottom in this guide so all commands work seemlessly.

## Prerequisites

The commands in this guide are known to work with the following prerequisites: 

### Resources

* CPU's: 2
* Memory: 4Gb

### Dependencies
* Ubuntu, Version 18.04
* Docker Community Editition, Version 20.10.12
* [Compose V2]((https://docs.docker.com/compose/cli-command/), Version 2.2.2

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

$ # Run all end-to-end tests
$ npm run e2e

$ # Or open the end-to-end test suite interactively
$ npm run e2e:open
```

### Developing the Frontend

You can switch to `ng serve` for better development experience.

```
$ # Reuse the environment.json file from the still running frontend container
$ curl http://localhost:4200/assets/environment.json > ./console/src/assets/environment.json

$ # Stop the frontend container
$ docker compose -f ./build/local/docker-compose-local.yml --profile frontend stop

$ # Change to the console directory
$ cd ./console

$ # Run the local server
$ ng serve

$ # Run all end-to-end tests
$ npm run e2e

$ # Or open the end-to-end test suite interactively
$ npm run e2e:open
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

### Mac M1 (Apple Silicon)

Bellow are some errors we faced with apple silicon.

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

### Remove the quickstart

```bash
$ docker compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend rm
```

If you are **confident** that you don't need to run the same ZITADEL instance again, go ahead and delete the `.keys` folder and reset the `environment.json` as well.

```bash
$ rm -rf .keys
```

```bash
$ git reset build/local/environment.json
```
