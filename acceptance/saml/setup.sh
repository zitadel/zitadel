#!/bin/sh

set -ex

PAT_FILE=${PAT_FILE:-../pat/zitadel-admin-sa.pat}
ZITADEL_API_URL="${ZITADEL_API_URL:-"http://localhost:8080"}"
LOGIN_URI="${LOGIN_URI:-"http://localhost:3000"}"
SAML_SP_METADATA="${SAML_SP_METADATA:-"http://samlsp:8081/saml/metadata"}"

if [ -z "${PAT}" ]; then
  echo "Reading PAT from file ${PAT_FILE}"
  PAT=$(cat ${PAT_FILE})
fi

#################################################################
# SAML Application
#################################################################

SAML_PROJECT_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_URL}/management/v1/projects" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{ \"name\": \"SAML\", \"projectRoleAssertion\": true, \"projectRoleCheck\": true, \"hasProjectCheck\": true, \"privateLabelingSetting\": \"PRIVATE_LABELING_SETTING_UNSPECIFIED\"}")
echo "Received SAML Project response: ${SAML_PROJECT_RESPONSE}"

SAML_PROJECT_ID=$(echo ${SAML_PROJECT_RESPONSE} | jq -r '. | .id')
echo "Received Project ID: ${SAML_PROJECT_ID}"

SAML_APP_RESPONSE=$(curl -s --request POST \
      --url "${ZITADEL_API_URL}/management/v1/projects/${SAML_PROJECT_ID}/apps/saml" \
      --header "Authorization: Bearer ${PAT}" \
      --header "Host: ${ZITADEL_API_DOMAIN}" \
      --header "Content-Type: application/json" \
      -d "{ \"name\": \"SAML\", \"metadataUrl\": \"${SAML_SP_METADATA}\", \"loginVersion\": { \"loginV2\": { \"baseUri\": \"${LOGIN_URI}\" }}}")
echo "Received SAML App response: ${SAML_APP_RESPONSE}"
