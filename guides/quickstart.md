# Quickstart with docker compose

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

ZITADEL is ready as soon as you can read

```text
++=========++
|| ZITADEL ||
|| STARTED ||
++=========++
```

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
&& docker-compose -f ./build/local/docker-compose-local.yml --profile database -p zitadel up --exit-code-from db-migrations \
&& sleep 5 \
&& docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend --profile setup -p zitadel up
```

For a more detailed guide take a look at the [development guide](./development.md)
