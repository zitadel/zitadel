#!/usr/bin/env bash

export USERPW=$(gopass show -o zitadel-secrets/zitadel/dev/moss@caos-demo.zitadel.dev)
export SERVICEACCOUNT_KEY=$(gopass show -o zitadel-secrets/zitadel/prod/e2e/cypressserviceuserkey)

npm install mochawesome  --save-dev

docker run \
--env CYPRESS_username="moss@caos-demo.zitadel.dev" \
--env CYPRESS_password="$USERPW" \
--env CYPRESS_consoleUrl="https://console.zitadel.dev" \
--env CYPRESS_projectName="newProject"  \
--env CYPRESS_serviceAccountKey="${SERVICEACCOUNT_KEY}" \
-it -v $PWD:/e2e -w /e2e cypress/included:8.0.0 
