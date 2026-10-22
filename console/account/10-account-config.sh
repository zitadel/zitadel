#!/bin/sh
set -eu
: "${ACCOUNT_ISSUER:?Set ACCOUNT_ISSUER}" "${ACCOUNT_API:?Set ACCOUNT_API}" "${ACCOUNT_CLIENT_ID:?Set ACCOUNT_CLIENT_ID}" "${ACCOUNT_LOGIN_ORG_ID:?Set ACCOUNT_LOGIN_ORG_ID}"
# These origins are also interpolated into CSP; accept origins only, never paths or header delimiters.
jq -en --arg issuer "$ACCOUNT_ISSUER" --arg api "$ACCOUNT_API" --arg clientid "$ACCOUNT_CLIENT_ID" --arg org "$ACCOUNT_LOGIN_ORG_ID" '
  def origin: test("^https://[A-Za-z0-9.-]+(:[0-9]+)?$");
  if ($issuer|origin) and ($api|origin) and ($org|test("^[0-9]+$")) and ($clientid|length > 0)
  then {issuer:$issuer,api:$api,clientid:$clientid,loginOrgId:$org} else error("Invalid Account environment") end
' > /usr/share/nginx/html/assets/environment.json
