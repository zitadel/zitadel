#!/bin/sh

set -ex

# The path to the instance admin service account
KEY=${KEY:-./machinekey/zitadel-admin-sa.json}
# The audience for which the key is used
AUDIENCE=${AUDIENCE:-http://localhost:8080}
# The Service can differ from the audience, for example in docker compose (http://zitadel:8080)
SERVICE=${SERVICE:-$AUDIENCE}

# Defer the Host header sent in requests that ZITADEL maps to an instance from the JWT audience
AUDIENCE_HOST="$(echo $AUDIENCE | cut -d/ -f3)"

# Create JWT from Admin SA KEY
JWT=$(zitadel-tools key2jwt --key ${KEY} --audience ${AUDIENCE})

# Get Token
TOKEN_RESPONSE=$(curl --request POST \
  --url ${SERVICE}/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header "Host: ${AUDIENCE_HOST}" \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile email' \
  --data assertion="${JWT}")

# Extract Token
TOKEN=$(echo -n ${TOKEN_RESPONSE} | jq -r '.access_token')

# Verify authentication
curl --request POST \
  --url ${SERVICE}/oidc/v1/userinfo \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}"

# Get default org
curl --request GET \
  --url ${SERVICE}/admin/v1/orgs/default \
  --header 'Accept: application/json' \
  --header "Authorization: Bearer ${TOKEN}" \
  --header "Host: ${AUDIENCE_HOST}"
