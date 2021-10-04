#!/usr/bin/env bash

set -a; source $2; set +a

npx cypress $1 --port 4201 --env org_owner_password="${E2E_ORG_OWNER_PW}",org_owner_viewer_password="${E2E_ORG_OWNER_VIEWER_PW}",org_project_creator_password="${E2E_ORG_PROJECT_CREATOR_PW}",consoleUrl=${E2E_CONSOLE_URL},apiCallsDomain="${E2E_API_CALLS_DOMAIN}",serviceAccountKey="${E2E_SERVICEACCOUNT_KEY}",zitadelProjectResourceId="${E2E_ZITADEL_PROJECT_RESOURCE_ID}"
