#!/bin/sh

if [ -f /.env-file/.env ]; then
    set -o allexport
    . /.env-file/.env
    set +o allexport
fi

if [ -n "${SYSTEM_USER_PRIVATE_KEY_FILE}" ]; then
    echo "SYSTEM_USER_PRIVATE_KEY_FILE=${SYSTEM_USER_PRIVATE_KEY_FILE} is set. Awaiting file and reading token."
    while [ ! -f "${SYSTEM_USER_PRIVATE_KEY_FILE}" ]; do
        sleep 2
    done
    echo "private key file found, reading private key"
    export SYSTEM_USER_PRIVATE_KEY=$(cat "${SYSTEM_USER_PRIVATE_KEY_FILE}")
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
