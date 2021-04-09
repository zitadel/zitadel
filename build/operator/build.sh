#!/bin/bash

set -e

./build/operator/prebuild.sh "./migrations"

CGO_ENABLED=0 go build \
  -a \
  -installsuffix cgo \
  -ldflags "$(./build/operator/ldflags.sh "${1}")" \
  -o zitadelctl \
  ./cmd/zitadelctl/main.go
