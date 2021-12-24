# Quickstart with docker compose

## Prerequisites

The only prerequisite you need to have installed is docker. The resource limit must at least be:

* CPU's: 2
* Memory: 4Gb

## Start ZITADEL

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

You can connect to [ZITADEL on localhost:4200](http://localhost:4200) after the frontend compiled  successfully. Initially it takes several minutes to start all containers.

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker compose -f ./build/local/docker-compose-local.yml --profile backend --profile frontend up
```

For a more detailed guide take a look at the [development guide](./development.md)

## FAQ

### Mac M1

Bellow are some error's we faced with apple silicon.

#### database-migrations don't start or stop without exit code

The problem is that the database has an error. You can simply restart the database with the following command:

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker compose -f ./build/local/docker-compose-local.yml restart db
```

#### API call's block and don't return any response

The problem is that the database has an error. You can simply restart the database with the following command:

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker compose -f ./build/local/docker-compose-local.yml restart db
```

### Build Errors

If you experience strange docker error you might need to check that `buildkit` is enabled.

Make sure to enable `"features": { "buildkit": true }` in your docker settings!

### Remove the quickstart

```Bash
docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend rm
```

If you are **confident** that you don't need to run the same ZITADEL instance again, go ahead and delete the `.keys` folder and reset the `environment.json` as well.

```Bash
rm -rf .keys
```

```Bash
git reset build/local/environment.json
```
