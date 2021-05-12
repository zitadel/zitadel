# Startup ZITADEL

## Prerequesits

* docker
* go

`go get github.com/rakyll/statik`

## ENV variables
```bash
#tracing is disabled locally
export ZITADEL_TRACING_TYPE=none
#metrics is disabled locally
export ZITADEL_METRICS_TYPE=none

#recommended log level for local is debug
export ZITADEL_LOG_LEVEL="debug"

#database connection (local cockroach insecure)
export ZITADEL_EVENTSTORE_HOST="localhost"
export ZITADEL_EVENTSTORE_PORT="26257"
export CR_SSL_MODE="disable"

#keys for security
export ZITADEL_KEY_PATH=".keys/local_keys.yaml"
export ZITADEL_USER_VERIFICATION_KEY="userverificationkey_1"
export ZITADEL_OTP_VERIFICATION_KEY="OTPVerificationKey_1"
export ZITADEL_OIDC_KEYS_ID="oidckey_1"
export ZITADEL_COOKIE_KEY="cookiekey_1"
export ZITADEL_CSRF_KEY="cookiekey_1"
export ZITADEL_IDP_CONFIG_VERIFICATION_KEY="idpconfigverificationkey_1"
export ZITADEL_DOMAIN_VERIFICATION_KEY="domainverificationkey_1"

#debug mode is used for notifications
export DEBUG_MODE=true
export CAOS_OIDC_DEV=true
export ZITADEL_CSRF_DEV=true

#TODO: needed for local?
export TWILIO_SENDER_NAME="CAOS AG"
export SMTP_HOST="smtp.gmail.com:465"
export SMTP_USER="zitadel@caos.ch"
export EMAIL_SENDER_ADDRESS="noreply@caos.ch"
export EMAIL_SENDER_NAME="CAOS AG"
export SMTP_TLS=true

#configuration for api/browser calls
export ZITADEL_DEFAULT_DOMAIN="zitadel.ch"
export ZITADEL_ISSUER="http://localhost:50002/oauth/v2/"
export ZITADEL_ACCOUNTS="http://localhost:50003/login"
export ZITADEL_AUTHORIZE="http://localhost:50002/oauth/v2"
export ZITADEL_OAUTH="http://localhost:50002/oauth/v2"
export ZITADEL_CONSOLE="http://localhost:4200"
export ZITADEL_COOKIE_DOMAIN="localhost"

#caching is used in UI's and API's
export ZITADEL_CACHE_MAXAGE=12h
export ZITADEL_CACHE_SHARED_MAXAGE=168h
export ZITADEL_SHORT_CACHE_MAXAGE=5m
export ZITADEL_SHORT_CACHE_SHARED_MAXAGE=15m

#console authorization configuration
export ZITADEL_CONSOLE_RESPONSE_TYPE="CODE"
export ZITADEL_CONSOLE_GRANT_TYPE="AUTHORIZATION_CODE"

export ZITADEL_CONSOLE_DEV_MODE=true
export ZITADEL_CONSOLE_ENV_DIR="${workspaceFolder}/console/src/assets/"
```

## Pre steps

### generate code

1. `go generate internal/statik/generate.go`
2. `go generate openapi/statik/generate.go`
3. `go generate internal/ui/login/statik/generate.go`
4. `go generate internal/notification/statik/generate.go`
5. `DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --target go-copy -o .`

### start database (cockroach)

```bash
docker run -d \                                  
--name=zitadel-db \
--hostname=zitadel-db \
-p 26257:26257 -p 8080:8080 \
-v "${GOPATH}/src/github.com/caos/zitadel/cockroach-data/citadel1:/cockroach/cockroach-data"  \
-v "${GOPATH}/src/github.com/caos/zitadel/.backups:/backups" \
cockroachdb/cockroach:v20.2.9 start-single-node --insecure --listen-addr=0.0.0.0
```

### execute migrations

`go generate migrations/cockroach/migrate_local.go`

