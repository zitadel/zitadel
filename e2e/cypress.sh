#!/usr/bin/env bash

ACTION=$1
ENVFILE=$2

shift
shift

projectRoot=".."

set -a; source $ENVFILE; set +a

NPX=""
if ! command -v cypress &> /dev/null; then
    NPX="npx"
fi

$NPX cypress $ACTION \
    --port ${E2E_CYPRESS_PORT} \
    --env org="${E2E_ORG}",org_owner_password="${E2E_ORG_OWNER_PW}",org_owner_viewer_password="${E2E_ORG_OWNER_VIEWER_PW}",org_project_creator_password="${E2E_ORG_PROJECT_CREATOR_PW}",login_policy_user_password="${E2E_LOGIN_POLICY_USER_PW}",password_complexity_user_password="${E2E_PASSWORD_COMPLEXITY_USER_PW}",consoleUrl=${E2E_CONSOLE_URL},apiUrl="${E2E_API_URL}",accountsUrl="${E2E_ACCOUNTS_URL}",issuerUrl="${E2E_ISSUER_URL}",serviceAccountKey="${E2E_SERVICEACCOUNT_KEY}",serviceAccountKeyPath="${E2E_SERVICEACCOUNT_KEY_PATH}",otherZitadelIdpInstance="${E2E_OTHER_ZITADEL_IDP_INSTANCE}",zitadelProjectResourceId="${E2E_ZITADEL_PROJECT_RESOURCE_ID}" \
    "$@"
