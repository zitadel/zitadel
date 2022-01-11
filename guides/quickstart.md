# Quickstart with docker compose

## Prerequisites

The only prerequisite you need fullfill, is that you need to have docker installed with support for compose and buildkit. The resource limit must at least be:

* CPU's: 2
* Memory: 4Gb

## Start ZITADEL

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

You can connect to [ZITADEL on localhost:4200](http://localhost:4200) after the frontend compiled  successfully. Initially it takes several minutes to start all containers.

<a name="compose-services"></a>
```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend up
```

## Developing ZITADEL

Instead of the profiles backend and frondend as described [above](#compose-services), use the profile e2e and detach from containers.

<a name="compose-e2e"></a>
```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml --profile e2e up -d
```

This also initializes data needed by Cypress end-to-end tests. Launch the Cypress test suite from the console directory:

```bash
cd ./console
./cypress.sh open local_local.env
```

You can run any test files except init.ts, as this is already run by the docker compose command shown [above](#compose-e2e) and only passes once.

Make changes to a service as you wish and rebuild and deploy the service using the following command from the project root directory:
```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml up -d --no-deps --build <compose service>
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
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml restart db
```

#### API call's block and don't return any response

The problem is that the database has a connection issues. You can simply restart the database with the following command:

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f ./build/local/docker-compose-local.yml restart db
```

### Build Errors

If you experience strange docker error you might need to check that `buildkit` is enabled.

Make sure to enable `"features": { "buildkit": true }` in your docker settings!

### Remove the quickstart

```Bash
docker compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend rm
```

If you are **confident** that you don't need to run the same ZITADEL instance again, go ahead and delete the `.keys` folder and reset the `environment.json` as well.

```Bash
rm -rf .keys
```

```Bash
git reset build/local/environment.json
```
