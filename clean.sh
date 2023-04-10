#!/bin/bash

set -eux

# Clean folders & generated files
rm -rf .artifacts/grpc
rm -rf console/src/app/proto/generated
rm -rf openapi/v2/zitadel
rm -rf docs/apis/assets
# rm -rf openapi/statik/statik.go
find pkg/grpc -name \*.pb.go -type f -delete
find pkg/grpc -name \*.pb.validate.go -type f -delete
find pkg/grpc -name \*.pb.authoptions.go -type f -delete
find pkg/grpc -name \*.pb.gw.go -type f -delete

# Create folders where needed
mkdir -p openapi/v2/zitadel
mkdir -p docs/apis/assets