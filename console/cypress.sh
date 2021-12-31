#!/usr/bin/env bash

ACTION=$1
ENVFILE=$2

shift
shift

set -a; source $ENVFILE; set +a

NPX=""
if ! command -v cypress &> /dev/null; then
    NPX="npx" 
fi

$NPX cypress $ACTION --port 4201 --env org_owner_password="${E2E_ORG_OWNER_PW}",org_owner_viewer_password="${E2E_ORG_OWNER_VIEWER_PW}",org_project_creator_password="${E2E_ORG_PROJECT_CREATOR_PW}",consoleUrl=${E2E_CONSOLE_URL},apiCallsDomain="${E2E_API_CALLS_DOMAIN}",serviceAccountKey="${E2E_SERVICEACCOUNT_KEY}",zitadelProjectResourceId="${E2E_ZITADEL_PROJECT_RESOURCE_ID}" "$@"
