#!/bin/sh
set -o allexport
. /.env-file/.env
set +o allexport

if [ -n "${ZITADEL_SERVICE_USER_TOKEN_PATH}" ] && [ -f "${ZITADEL_SERVICE_USER_TOKEN_PATH}" ]; then
  export ZITADEL_SERVICE_USER_TOKEN=$(cat "${ZITADEL_SERVICE_USER_TOKEN_PATH}")
fi

exec node apps/login/server.js
