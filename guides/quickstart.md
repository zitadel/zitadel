# Quickstart with docker compose

You can start ZITADEL with a simple docker compose up.

The services are configured to restart if an error occurs.

In the following script the basic setup of the database is executed before ZITADEL starts. Execute the statement from the root of ZITADEL.

You can connect to [ZITADEL on localhost:4200](http://localhost:4200) as soon as the following text appears:

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

## FAQ

### Build Errors

If you experience strange docker error you might need to check that `buildkit` is enabled.

Make sure to enable `"features": { "buildkit": true }` in your docker settings!

### Error 412

If the line `backend (412), environment (000) or console (000) not ready yet ==> retrying in 5 seconds` contains a `412` you struck a race condition within the start process (It should exit 1 when this happens). This will be fixed in a future release. You can work around it by restarting the containers.

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend --profile setup -p zitadel up
```

### Remove the quickstart

```Bash
docker-compose -f ./build/local/docker-compose-local.yml --profile database --profile init-backend --profile init-frontend --profile backend --profile frontend --profile setup -p zitadel rm
docker system prune -a

If you are **confident** that you don't need to run the same ZITADEL instance again, go ahead and delete the `.keys` folder and reset the `environment.json` as well.

```Bash
rm -rf .keys
```

```Bash
git reset build/local/environment.json
```
