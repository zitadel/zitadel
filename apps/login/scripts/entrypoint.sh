#!/bin/sh

if [ -f /.env-file/.env ]; then
    set -o allexport
    . /.env-file/.env
    set +o allexport
fi

if [ -n "${ZITADEL_SERVICE_USER_TOKEN_FILE}" ]; then
    echo "ZITADEL_SERVICE_USER_TOKEN_FILE=${ZITADEL_SERVICE_USER_TOKEN_FILE} is set. Awaiting file and reading token."
    while [ ! -f "${ZITADEL_SERVICE_USER_TOKEN_FILE}" ]; do
        sleep 2
    done
    echo "token file found, reading token"
    export ZITADEL_SERVICE_USER_TOKEN=$(cat "${ZITADEL_SERVICE_USER_TOKEN_FILE}")
fi

exec $@
