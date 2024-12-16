#!/bin/sh

set -ex

PAT_FILE=${PAT_FILE:-./pat/zitadel-admin-sa.pat}
ZITADEL_API_PROTOCOL="${ZITADEL_API_PROTOCOL:-http}"
ZITADEL_API_DOMAIN="${ZITADEL_API_DOMAIN:-localhost}"
ZITADEL_API_PORT="${ZITADEL_API_PORT:-8080}"
ZITADEL_API_URL="${ZITADEL_API_URL:-${ZITADEL_API_PROTOCOL}://${ZITADEL_API_DOMAIN}:${ZITADEL_API_PORT}}"
ZITADEL_API_INTERNAL_URL="${ZITADEL_API_INTERNAL_URL:-${ZITADEL_API_URL}}"
SINK_EMAIL_INTERNAL_URL="${SINK_EMAIL_INTERNAL_URL:-"http://sink:3333/email"}"
SINK_SMS_INTERNAL_URL="${SINK_SMS_INTERNAL_URL:-"http://sink:3333/sms"}"
SINK_NOTIFICATION_URL="${SINK_NOTIFICATION_URL:-"http://localhost:3333/notification"}"

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

#################################################################
# Environment files
#################################################################

WRITE_ENVIRONMENT_FILE=${WRITE_ENVIRONMENT_FILE:-$(dirname "$0")/../apps/login/.env.local}
echo "Writing environment file to ${WRITE_ENVIRONMENT_FILE} when done."
WRITE_TEST_ENVIRONMENT_FILE=${WRITE_TEST_ENVIRONMENT_FILE:-$(dirname "$0")/../acceptance/tests/.env.local}
echo "Writing environment file to ${WRITE_TEST_ENVIRONMENT_FILE} when done."

echo "ZITADEL_API_URL=${ZITADEL_API_URL}
ZITADEL_SERVICE_USER_ID=${ZITADEL_SERVICE_USER_ID}
ZITADEL_SERVICE_USER_TOKEN=${PAT}
SINK_NOTIFICATION_URL=${SINK_NOTIFICATION_URL}
DEBUG=true"| tee "${WRITE_ENVIRONMENT_FILE}" "${WRITE_TEST_ENVIRONMENT_FILE}" > /dev/null
echo "Wrote environment file ${WRITE_ENVIRONMENT_FILE}"
cat ${WRITE_ENVIRONMENT_FILE}

echo "Wrote environment file ${WRITE_TEST_ENVIRONMENT_FILE}"
cat ${WRITE_TEST_ENVIRONMENT_FILE}

#################################################################
# SMS provider with HTTP
#################################################################

SMSHTTP_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/admin/v1/sms/http" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"endpoint\": \"${SINK_SMS_INTERNAL_URL}\", \"description\": \"test\"}")
echo "Received SMS HTTP response: ${SMSHTTP_RESPONSE}"

SMSHTTP_ID=$(echo ${SMSHTTP_RESPONSE} | jq -r '. | .id')
echo "Received SMS HTTP ID: ${SMSHTTP_ID}"

SMS_ACTIVE_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/admin/v1/sms/${SMSHTTP_ID}/_activate" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json")
echo "Received SMS active response: ${SMS_ACTIVE_RESPONSE}"

#################################################################
# Email provider with HTTP
#################################################################

EMAILHTTP_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/admin/v1/email/http" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"endpoint\": \"${SINK_EMAIL_INTERNAL_URL}\", \"description\": \"test\"}")
echo "Received Email HTTP response: ${EMAILHTTP_RESPONSE}"

EMAILHTTP_ID=$(echo ${EMAILHTTP_RESPONSE} | jq -r '. | .id')
echo "Received Email HTTP ID: ${EMAILHTTP_ID}"

EMAIL_ACTIVE_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/admin/v1/email/${EMAILHTTP_ID}/_activate" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json")
echo "Received Email active response: ${EMAIL_ACTIVE_RESPONSE}"

#################################################################
# Wait for projection of default organization in ZITADEL
#################################################################

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

