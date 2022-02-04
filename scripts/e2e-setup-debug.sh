#!/bin/bash

set -ex

set -a
source ./build/local/local.env
source ./console/e2e.env
set +a

dlv debug --api-version 2 --headless --listen 127.0.0.1:2345 ./cmd/e2e-setup/*.go -- --setup-files "./cmd/zitadel/authz.yaml" --setup-files "./cmd/zitadel/system-defaults.yaml" --setup-files "./cmd/zitadel/setup.yaml" --setup-files "./cmd/e2e-setup/e2e.yaml"
