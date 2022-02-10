#!/bin/bash

set -e

projectRoot="."

set -a
source ./build/local/local.env
source <(grep -v E2E_ZITADEL_PROJECT_RESOURCE_ID ./console/e2e.env)
source ./build/local/e2e-setup-overwrite.env
set +a

go run ./cmd/e2e-setup/*.go --setup-files "./cmd/zitadel/authz.yaml" --setup-files "./cmd/zitadel/system-defaults.yaml" --setup-files "./cmd/zitadel/setup.yaml" --setup-files "./cmd/e2e-setup/e2e.yaml"
