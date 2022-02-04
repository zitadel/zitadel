#!/bin/bash

set -a
source ./build/local/local.env
source ./build/local/e2e.env
set +a

go run ./cmd/e2e-setup/* --setup-files "./cmd/zitadel/authz.yaml" --setup-files "./cmd/zitadel/system-defaults.yaml" --setup-files "./cmd/zitadel/setup.yaml" --setup-files "./cmd/zitadel/startup.yaml"
