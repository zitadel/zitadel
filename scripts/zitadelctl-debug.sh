#!/bin/bash

set -e

./build/operator/prebuild.sh "./migrations"

dlv debug --api-version 2 --headless --listen 127.0.0.1:2345 ./cmd/zitadelctl -- "$@"
