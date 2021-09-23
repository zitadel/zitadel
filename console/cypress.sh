#!/usr/bin/env bash

set -a; source $2; set +a

npx cypress $1 --port 4201 --env org_owner_password="${ORG_OWNER_PW}",org_owner_viewer_password="${ORG_OWNER_VIEWER_PW}",org_project_creator_password="${ORG_PROJECT_CREATOR_PW}",consoleUrl=${CONSOLE_URL},apiCallsDomain="${API_CALLS_DOMAIN}",projectName="${PROJECT_NAME}",serviceAccountKey="${SERVICEACCOUNT_KEY}",zitadelProjectResourceId="${ZITADEL_PROJECT_RESOURCE_ID}"
