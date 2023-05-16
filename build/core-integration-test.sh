#!/bin/bash

case ${DB_FLAVOR} in
    cockroach)
        cockroach start-single-node --insecure --listen-addr=localhost:36257 --sql-addr=localhost:26257 --http-addr=localhost:9090 --background
    ;;
    postgres)
        (exec /docker-entrypoint.sh postgres &> /dev/null &)
    ;;
esac

go build -o zitadel main.go

./zitadel init --config internal/integration/config/zitadel.yaml --config internal/integration/config/${DB_FLAVOR}.yaml
./zitadel setup --masterkeyFromEnv --config internal/integration/config/zitadel.yaml --config internal/integration/config/${DB_FLAVOR}.yaml

go test -tags=integration -race -parallel 1 -v -coverprofile=profile.cov -coverpkg=./... ./internal/integration ./internal/api/grpc/...