#!/bin/sh
set -o allexport
. /.env-file/.env
set +o allexport

if [ -n "${ZITADEL_SERVICE_USER_TOKEN_FILE}" ] && [ -f "${ZITADEL_SERVICE_USER_TOKEN_FILE}" ]; then
  echo "ZITADEL_SERVICE_USER_TOKEN_FILE=${ZITADEL_SERVICE_USER_TOKEN_FILE} is set and file exists, setting ZITADEL_SERVICE_USER_TOKEN to the files content"
  export ZITADEL_SERVICE_USER_TOKEN=$(cat "${ZITADEL_SERVICE_USER_TOKEN_FILE}")
fi

exec node /runtime/server.js
