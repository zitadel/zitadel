# Quickstart with docker compose

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

In order to use Docker and Compose with buildkit enabled, export two environment variables for your current shell:

```bash
$ export DOCKER_BUILDKIT=1 
$ export COMPOSE_DOCKER_CLI_BUILD=1
```

## Starting ZITADEL

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

You can connect to [ZITADEL on localhost:4200](http://localhost:4203) after the frontend compiled successfully. Initially it takes several minutes to start all containers.

<a name="compose-services"></a>
```bash
$ docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend up --detach
```

## Developing ZITADEL

If you want to make changes to ZITADEL, we recommend running the end-to-end tests against it. 

### Test Prerequisites

Additionally to the prerequsites described [above](#prerequisites), the end-to-end tests are know to work with the following dependencies:

* NodeJS, Version 14.17.6
* NPM, Version 6.14.15

### Running End-to-End Tests

The following command creates all objects in your local ZITADEL needed to run the end-to-end tests. The command may take several minutes because it awaits eventual consistency but you only have to execute it once. 

```bash
$ ./scripts/e2e-setup.sh
```

Now you are ready to actually run the test suite.

```bash
$ # Change directory to ./console
$ cd ./console

$ # Install dev dependencies if you haven't done so already
$ npm install --only development

$ # Run all end-to-end tests
$ npm run e2e

$ # Or open the end-to-end test suite interactively
$ npm run e2e:open
```

### Redeploying a Service

Make changes to a service as you wish and rebuild and deploy it using the following command from the project root directory:

```bash
$ docker compose -f ./build/local/docker-compose-local.yml up -d --no-deps --build <compose service>
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
