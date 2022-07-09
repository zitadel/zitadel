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
    --port "${E2E_CYPRESSPORT}" \
    --env org="${ZITADEL_E2E_ORG}",org_owner_password="${ZITADEL_E2E_ORGOWNERPW}",org_owner_viewer_password="${ZITADEL_E2E_ORGOWNERVIEWERPW}",org_project_creator_password="${ZITADEL_E2E_ORGPROJECTCREATORPW}",login_policy_user_password="${ZITADEL_E2E_LOGINPOLICYUSERPW}",password_complexity_user_password="${ZITADEL_E2E_PASSWORDCOMPLEXITYUSERPW}",consoleUrl="${ZITADEL_E2E_CONSOLEURL}",apiUrl="${ZITADEL_E2E_APIURL}",accountsUrl="${ZITADEL_E2E_ACCOUNTSURL}",issuerUrl="${ZITADEL_E2E_ISSUERURL}",serviceAccountKeyPath="${ZITADEL_E2E_MACHINEKEYPATH}",otherZitadelIdpInstance="${ZITADEL_E2E_OTHERZITADELIDPINSTANCE}",zitadelProjectResourceId="${ZITADEL_E2E_ZITADELPROJECTRESOURCEID}" \
    "$@"
