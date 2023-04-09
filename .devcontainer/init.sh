#!/bin/bash

set -eux

# Clean folders & generated files
rm -rf .artifacts/grpc
rm -rf console/src/app/proto/generated
rm -rf openapi/v2/zitadel
rm -rf docs/apis/assets
rm -rf openapi/statik/statik.go
find pkg/grpc -name \*.pb.go -type f -delete
find pkg/grpc -name \*.pb.validate.go -type f -delete
find pkg/grpc -name \*.pb.authoptions.go -type f -delete
find pkg/grpc -name \*.pb.gw.go -type f -delete

# Create folders where needed
mkdir -p openapi/v2/zitadel
mkdir -p docs/apis/assets

# Install Node deps & generate gRPC client
(cd console ; npm ci)
(cd console ; npm run generate)

# Install Go deps & generate gRPC server
go mod download
go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
go install github.com/rakyll/statik@v0.1.7

for i in $(find proto/zitadel -iname *.proto); do export PATH=$(go env GOPATH)/bin:$PATH && ./node_modules/.bin/buf generate ${i}; done
cp .artifacts/grpc/zitadel/auth.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/auth
cp .artifacts/grpc/zitadel/admin.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/admin
cp .artifacts/grpc/zitadel/management.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/management
cp .artifacts/grpc/zitadel/system.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/system
cp -R .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc

go generate internal/api/ui/login/statik/generate.go
go generate internal/api/ui/login/static/generate.go
go generate internal/notification/statik/generate.go
go generate internal/statik/generate.go

cp .artifacts/grpc/zitadel/*.swagger.json openapi/v2/zitadel
go generate openapi/statik/generate.go
go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

exit 0