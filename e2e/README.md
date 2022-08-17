# Run  e2e Tests

```bash
docker compose run e2e
```

# Cleanup e2e Tests

```bash
docker compose down
```


# forward Cypress GUI on Mac
Install Xquarts
```bash
brew cask install xquartz 
```
configure X11 preferences to "allow connections from network clients"
Install XQuartz as per https://sourabhbajaj.com/blog/2017/02/07/gui-applications-docker-mac/

set IP and DISPLAY variable and alow xhost communication

```bash
IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
DISPLAY=$IP:0

xhost + $IP
```

start cypress with xforward
```bash
COMPOSE_DOCKER_CLI_BUILD=1 docker compose -f docker-compose.yaml  -f docker-compose-cypress-open.yaml up
```
