#!/bin/bash

set -e

CGO_ENABLED=0 go build \
  -a \
  -installsuffix cgo \
  -ldflags "$(./build/operator/ldflags.sh "${1}")" \
  -o zitadelctl \
  ./cmd/zitadelctl/main.go
