## start with compose.env for automated tests
COMPOSE_DOCKER_CLI_BUILD=1 docker compose -f docker-compose.yaml -f docker-compose-overwrite.yaml --env-file compose.env run e2e


## cleanup containers and volumes
docker rm -f $(docker ps -a -q)
docker volume rm $(docker volume ls -q)
