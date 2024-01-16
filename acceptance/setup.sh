#!/bin/sh

set -e

KEY=${KEY:-./machinekey/zitadel-admin-sa.json}
echo "Using key path ${KEY} to the instance admin service account."

AUDIENCE=${AUDIENCE:-http://localhost:8080}
echo "Using audience ${AUDIENCE} for which the key is used."

SERVICE=${SERVICE:-$AUDIENCE}
echo "Using the service ${SERVICE} to connect to ZITADEL. For example in docker compose this can differ from the audience."

WRITE_ENVIRONMENT_FILE=${WRITE_ENVIRONMENT_FILE:-$(dirname "$0")/../apps/login/.env.acceptance}
echo "Writing environment file to ${WRITE_ENVIRONMENT_FILE} when done."

AUDIENCE_HOST="$(echo $AUDIENCE | cut -d/ -f3)"
echo "Deferred the Host header ${AUDIENCE_HOST} which will be sent in requests that ZITADEL then maps to a virtual instance"

JWT=$(zitadel-tools key2jwt --key ${KEY} --audience ${AUDIENCE})
echo "Created JWT from Admin service account key ${JWT}"

TOKEN_RESPONSE=$(curl -s --request POST \
  --url ${SERVICE}/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header "Host: ${AUDIENCE_HOST}" \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud' \
  --data assertion="${JWT}")
echo "Got response from token endpoint:"
echo "${TOKEN_RESPONSE}" | jq

TOKEN=$(echo -n ${TOKEN_RESPONSE} | jq --raw-output '.access_token')
echo "Extracted access token ${TOKEN}"

ORG_RESPONSE=$(curl -s --request GET \
  --url ${SERVICE}/admin/v1/orgs/default \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}")
echo "Got default org response:"
echo "${ORG_RESPONSE}" | jq

ORG_ID=$(echo -n ${ORG_RESPONSE} | jq --raw-output '.org.id')
echo "Extracted default org id ${ORG_ID}"

echo "ZITADEL_API_URL=${AUDIENCE}
ZITADEL_ORG_ID=${ORG_ID}
ZITADEL_SERVICE_USER_TOKEN=${TOKEN}" > ${WRITE_ENVIRONMENT_FILE}
echo "Wrote environment file ${WRITE_ENVIRONMENT_FILE}"
cat ${WRITE_ENVIRONMENT_FILE}

if ! grep -q 'localhost' ${WRITE_ENVIRONMENT_FILE}; then
  echo "Not developing against localhost, so creating a human user might not be necessary"
  exit 0
fi

HUMAN_USER_USERNAME="zitadel-admin@zitadel.localhost"
HUMAN_USER_PASSWORD="Password1!"

HUMAN_USER_PAYLOAD=$(cat << EOM
{
  "userName": "${HUMAN_USER_USERNAME}",
  "profile": {
    "firstName": "ZITADEL",
    "lastName": "Admin",
    "displayName": "ZITADEL Admin",
    "preferredLanguage": "en"
  },
  "email": {
    "email": "zitadel-admin@zitadel.localhost",
    "isEmailVerified": true
  },
  "password": "${HUMAN_USER_PASSWORD}",
  "passwordChangeRequired": false
}
EOM
)
echo "Creating human user"
echo "${HUMAN_USER_PAYLOAD}" | jq

HUMAN_USER_RESPONSE=$(curl -s --request POST \
  --url ${SERVICE}/management/v1/users/human/_import \
  --header 'Content-Type: application/json' \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}" \
  --data-raw "${HUMAN_USER_PAYLOAD}")
echo "Create human user response"
echo "${HUMAN_USER_RESPONSE}" | jq

if [ "$(echo -n "${HUMAN_USER_RESPONSE}" | jq --raw-output '.code')" == "6" ]; then
  echo "admin user already exists"
  exit 0
fi

HUMAN_USER_ID=$(echo -n ${HUMAN_USER_RESPONSE} | jq --raw-output '.userId')
echo "Extracted human user id ${HUMAN_USER_ID}"

HUMAN_ADMIN_PAYLOAD=$(cat << EOM
{
  "userId": "${HUMAN_USER_ID}",
  "roles": [
    "IAM_OWNER"
  ]
}
EOM
)
echo "Granting iam owner to human user"
echo "${HUMAN_ADMIN_PAYLOAD}" | jq

HUMAN_ADMIN_RESPONSE=$(curl -s --request POST \
  --url ${SERVICE}/admin/v1/members \
  --header 'Content-Type: application/json' \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}" \
  --data-raw "${HUMAN_ADMIN_PAYLOAD}")

echo "Grant iam owner to human user response"
echo "${HUMAN_ADMIN_RESPONSE}" | jq

echo "You can now log in at ${AUDIENCE}/ui/login"
echo "username: ${HUMAN_USER_USERNAME}"
echo "password: ${HUMAN_USER_PASSWORD}"