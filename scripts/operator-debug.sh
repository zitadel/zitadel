#!/bin/bash

set -e

go build -gcflags "all=-N -l" -o /tmp/zitadeloperator-debug ./cmd/operator-debug
dlv exec --api-version 2 --headless --listen 127.0.0.1:2345 /tmp/zitadeloperator-debug -- "$@"
