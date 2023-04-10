#!/bin/bash

set -eux

go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
go install github.com/bufbuild/buf/cmd/buf@v1.14.0

# This a dirty workaround for our own auth generator problem
for i in $(find proto/zitadel -iname *.proto); do buf generate ${i}; done
cp .artifacts/grpc/zitadel/auth.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/auth
cp .artifacts/grpc/zitadel/admin.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/admin
cp .artifacts/grpc/zitadel/management.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/management
cp .artifacts/grpc/zitadel/system.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/system
mkdir -p pkg/grpc
cp -R .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc
