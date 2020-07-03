#! /bin/sh

set -eux

GEN_PATH=console/src/app/proto/generated

echo "Remove old files"
rm -rf $GEN_PATH

echo "Create folders"
mkdir -p $GEN_PATH

echo "Download additional protofiles"
wget -P tmp/validate https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.4.0/validate/validate.proto
wget -P tmp/protoc-gen-swagger/options https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/annotations.proto
wget -P tmp/protoc-gen-swagger/options https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/openapiv2.proto

echo "Generate grpc"

protoc \
  -I=/usr/local/include \
  -I=pkg/management/api/proto \
  -I=pkg/auth/api/proto \
  -I=pkg/admin/api/proto \
  -I=internal/protoc/protoc-gen-authoption \
  -I=console/node_modules/google-proto-files \
  -I=tmp \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  pkg/management/api/proto/*.proto \
  pkg/admin/api/proto/*.proto \
  pkg/auth/api/proto/*.proto