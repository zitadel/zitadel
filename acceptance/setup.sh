#!/bin/sh

set -ex

PAT_FILE=${PAT_FILE:-./pat/zitadel-admin-sa.pat}
ZITADEL_API_PROTOCOL="${ZITADEL_API_PROTOCOL:-http}"
ZITADEL_API_DOMAIN="${ZITADEL_API_DOMAIN:-localhost}"
ZITADEL_API_PORT="${ZITADEL_API_PORT:-8080}"
ZITADEL_API_URL="${ZITADEL_API_URL:-${ZITADEL_API_PROTOCOL}://${ZITADEL_API_DOMAIN}:${ZITADEL_API_PORT}}"
ZITADEL_API_INTERNAL_URL="${ZITADEL_API_INTERNAL_URL:-${ZITADEL_API_URL}}"

if [ -z "${PAT}" ]; then
  echo "Reading PAT from file ${PAT_FILE}"
  PAT=$(cat ${PAT_FILE})
fi

if [ -z "${ZITADEL_SERVICE_USER_ID}" ]; then
  echo "Reading ZITADEL_SERVICE_USER_ID from userinfo endpoint"
  USERINFO_RESPONSE=$(curl -s --request POST \
    --url "${ZITADEL_API_INTERNAL_URL}/oidc/v1/userinfo" \
    --header "Authorization: Bearer ${PAT}" \
    --header "Host: ${ZITADEL_API_DOMAIN}")
  echo "Received userinfo response: ${USERINFO_RESPONSE}"
  ZITADEL_SERVICE_USER_ID=$(echo "${USERINFO_RESPONSE}" | jq --raw-output '.sub')
fi

WRITE_ENVIRONMENT_FILE=${WRITE_ENVIRONMENT_FILE:-$(dirname "$0")/../apps/login/.env.local}
echo "Writing environment file to ${WRITE_ENVIRONMENT_FILE} when done."

echo "ZITADEL_API_URL=${ZITADEL_API_URL}
ZITADEL_SERVICE_USER_ID=${ZITADEL_SERVICE_USER_ID}
ZITADEL_SERVICE_USER_TOKEN=${PAT}
DEBUG=true" > ${WRITE_ENVIRONMENT_FILE}

echo "Wrote environment file ${WRITE_ENVIRONMENT_FILE}"
cat ${WRITE_ENVIRONMENT_FILE}

DEFAULTORG_RESPONSE_RESULTS=0
# waiting for default organization
until [ ${DEFAULTORG_RESPONSE_RESULTS} -eq 1 ]
do
  DEFAULTORG_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/v2/organizations/_search" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"queries\": [{\"defaultQuery\":{}}]}" )
  echo "Received default organization response: ${DEFAULTORG_RESPONSE}"
  DEFAULTORG_RESPONSE_RESULTS=$(echo $DEFAULTORG_RESPONSE | jq -r '.result | length')
  echo "Received default organization response result: ${DEFAULTORG_RESPONSE_RESULTS}"
done
