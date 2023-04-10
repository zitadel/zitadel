#!/bin/bash

set -eux

# Install Go deps & generate gRPC server
go mod download

# Go bindata is currently removed and the file template.gen.go is commited to git.
go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
go install github.com/rakyll/statik@v0.1.7
go install github.com/bufbuild/buf/cmd/buf@v1.14.0

# This a dirty workaround for our own auth generator problem
for i in $(find proto/zitadel -iname *.proto); do buf generate ${i}; done
cp .artifacts/grpc/zitadel/auth.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/auth
cp .artifacts/grpc/zitadel/admin.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/admin
cp .artifacts/grpc/zitadel/management.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/management
cp .artifacts/grpc/zitadel/system.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/system
mkdir -p pkg/grpc
cp -R .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc

# go generate internal/api/ui/login/statik/generate.go
# # go generate internal/api/ui/login/static/generate.go # This looks like it should be go generate internal/api/ui/login/static/resources/generate.go but i used this from prod
# go generate internal/notification/statik/generate.go
# go generate internal/statik/generate.go

# cp .artifacts/grpc/zitadel/*.swagger.json openapi/v2/zitadel
# go generate openapi/statik/generate.go
# go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

# exit 0