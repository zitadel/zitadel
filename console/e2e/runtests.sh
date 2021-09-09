#!/usr/bin/env bash

export USERPW=$(gopass show -o zitadel-secrets/zitadel/dev/moss@caos-demo.zitadel.dev)

npm install mochawesome  --save-dev

docker run \
--env CYPRESS_username="moss@caos-demo.zitadel.dev" \
--env CYPRESS_password="$USERPW" \
--env CYPRESS_consoleUrl="https://console.zitadel.dev" \
--env CYPRESS_projectName="newProject"  \
-it -v $PWD:/e2e -w /e2e cypress/included:8.0.0 
