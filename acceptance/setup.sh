#!/bin/sh

set -ex

KEY=${KEY:-./machinekey/zitadel-admin-sa.json}
echo "Using key path ${KEY} to the instance admin service account."

AUDIENCE=${AUDIENCE:-http://localhost:8080}
echo "Using audience ${AUDIENCE} for which the key is used."

SERVICE=${SERVICE:-$AUDIENCE}
echo "Using the service ${SERVICE} to connect to ZITADEL. For example in docker compose this can differ from the audience."

WRITE_ENVIRONMENT_FILE=${WRITE_ENVIRONMENT_FILE:-$(dirname "$0")/../apps/login/.env.local}
echo "Writing environment file to ${WRITE_ENVIRONMENT_FILE} when done."

AUDIENCE_HOST="$(echo $AUDIENCE | cut -d/ -f3)"
echo "Deferred the Host header ${AUDIENCE_HOST} which will be sent in requests that ZITADEL then maps to a virtual instance"

JWT=$(zitadel-tools key2jwt --key ${KEY} --audience ${AUDIENCE})
echo "Created JWT from Admin service account key ${JWT}"

TOKEN_RESPONSE=$(curl --request POST \
  --url ${SERVICE}/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header "Host: ${AUDIENCE_HOST}" \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud' \
  --data assertion="${JWT}")
echo "Got response from token endpoint:"
echo "${TOKEN_RESPONSE}" | jq

TOKEN=$(echo -n ${TOKEN_RESPONSE} | jq -r '.access_token')
echo "Extracted access token ${TOKEN}"

ORG_RESPONSE=$(curl --request GET \
  --url ${SERVICE}/admin/v1/orgs/default \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}")
echo "Got default org response:"
echo "${ORG_RESPONSE}" | jq

ORG_ID=$(echo -n ${ORG_RESPONSE} | jq -r '.org.id')
echo "Extracted default org id ${ORG_ID}"

echo "ZITADEL_API_URL=${AUDIENCE}
ZITADEL_ORG_ID=${ORG_ID}
ZITADEL_SERVICE_USER_TOKEN=${TOKEN}" > ${WRITE_ENVIRONMENT_FILE}
echo "Wrote environment file ${WRITE_ENVIRONMENT_FILE}"
cat ${WRITE_ENVIRONMENT_FILE}