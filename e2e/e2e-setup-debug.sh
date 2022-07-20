#!/bin/bash

set -ex

export projectRoot="."

set -a
source ./e2e/local.env
set +a

dlv debug --api-version 2 --headless --listen 127.0.0.1:2345 ./cmd/e2e-setup/*.go "$@"
