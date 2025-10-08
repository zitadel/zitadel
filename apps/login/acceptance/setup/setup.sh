#!/bin/sh

set -e pipefail

PAT_FILE=${PAT_FILE:-./pat/zitadel-admin-sa.pat}
LOGIN_BASE_URL=${LOGIN_BASE_URL:-"http://localhost:3000"}
ZITADEL_API_PROTOCOL="${ZITADEL_API_PROTOCOL:-http}"
ZITADEL_API_DOMAIN="${ZITADEL_API_DOMAIN:-localhost}"
ZITADEL_API_PORT="${ZITADEL_API_PORT:-8080}"
ZITADEL_API_URL="${ZITADEL_API_URL:-${ZITADEL_API_PROTOCOL}://${ZITADEL_API_DOMAIN}:${ZITADEL_API_PORT}}"
ZITADEL_API_INTERNAL_URL="${ZITADEL_API_INTERNAL_URL:-${ZITADEL_API_URL}}"
SINK_EMAIL_INTERNAL_URL="${SINK_EMAIL_INTERNAL_URL:-"http://sink:3333/email"}"
SINK_SMS_INTERNAL_URL="${SINK_SMS_INTERNAL_URL:-"http://sink:3333/sms"}"
SINK_NOTIFICATION_URL="${SINK_NOTIFICATION_URL:-"http://localhost:3333/notification"}"
WRITE_ENVIRONMENT_FILE=${WRITE_ENVIRONMENT_FILE:-$(dirname "$0")/../apps/login/.env.test.local}

if [ -z "${PAT}" ]; then
  echo "Reading PAT from file ${PAT_FILE}"
  PAT=$(cat ${PAT_FILE})
fi

#################################################################
# ServiceAccount as Login Client
#################################################################

SERVICEACCOUNT_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/management/v1/users/machine" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"userName\": \"login\",  \"name\": \"Login v2\",  \"description\": \"Serviceaccount for Login v2\", \"accessTokenType\": \"ACCESS_TOKEN_TYPE_BEARER\"}")
echo "Received ServiceAccount response: ${SERVICEACCOUNT_RESPONSE}"

SERVICEACCOUNT_ID=$(echo ${SERVICEACCOUNT_RESPONSE} | jq -r '. | .userId')
echo "Received ServiceAccount ID: ${SERVICEACCOUNT_ID}"

MEMBER_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/admin/v1/members" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"userId\": \"${SERVICEACCOUNT_ID}\",  \"roles\": [\"IAM_LOGIN_CLIENT\"]}")
echo "Received Member response: ${MEMBER_RESPONSE}"

SA_PAT_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_INTERNAL_URL}/management/v1/users/${SERVICEACCOUNT_ID}/pats" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{\"expirationDate\": \"2519-04-01T08:45:00.000000Z\"}")
echo "Received Member response: ${MEMBER_RESPONSE}"

SA_PAT=$(echo ${SA_PAT_RESPONSE} | jq -r '. | .token')
echo "Received ServiceAccount Token: ${SA_PAT}"

#################################################################
# Environment files
#################################################################

echo "Writing environment file ${WRITE_ENVIRONMENT_FILE}."

echo "ZITADEL_API_URL=${ZITADEL_API_URL}
ZITADEL_SERVICE_USER_TOKEN=${SA_PAT}
ZITADEL_ADMIN_TOKEN=${PAT}
SINK_NOTIFICATION_URL=${SINK_NOTIFICATION_URL}
EMAIL_VERIFICATION=true
DEBUG=false
LOGIN_BASE_URL=${LOGIN_BASE_URL}
NODE_TLS_REJECT_UNAUTHORIZED=0
ZITADEL_ADMIN_USER=${ZITADEL_ADMIN_USER:-"zitadel-admin@zitadel.localhost"}
NEXT_PUBLIC_BASE_PATH=/ui/v2/login
" > ${WRITE_ENVIRONMENT_FILE}

echo "Wrote environment file ${WRITE_ENVIRONMENT_FILE}"
cat ${WRITE_ENVIRONMENT_FILE}

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

