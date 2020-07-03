#! /bin/sh

set -eux

GOPATH=${GOPATH:-~/go}

GEN_PATH=console/src/app/proto/generated

echo "Remove old files"
rm -rf $GEN_PATH

echo "Create folders"
mkdir -p $GEN_PATH

echo "Generate grpc"

protoc \
  -I=/usr/local/include \
  -I=pkg/management/api/proto \
  -I=pkg/auth/api/proto \
  -I=pkg/admin/api/proto \
  -I=internal/protoc/protoc-gen-authoption \
  -I=console/node_modules/google-proto-files \
  -I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.4.0 \
  -I=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6 \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  pkg/management/api/proto/*.proto \
  pkg/admin/api/proto/*.proto \
  pkg/auth/api/proto/*.proto

echo "Generate annotations js file (compatibility)"

mkdir -p $GEN_PATH/google/api/
touch $GEN_PATH/google/api/annotations_pb.js
echo "export {}" > $GEN_PATH/google/api/annotations_pb.d.ts

mkdir -p $GEN_PATH/validate
touch $GEN_PATH/validate/validate_pb.js
echo "export {}" > $GEN_PATH/validate/validate_pb.d.ts

mkdir -p $GEN_PATH/protoc-gen-swagger/options
touch $GEN_PATH/protoc-gen-swagger/options/annotations_pb.js
echo "export {}" > $GEN_PATH/protoc-gen-swagger/options/annotations_pb.d.ts

mkdir -p $GEN_PATH/authoption
touch $GEN_PATH/authoption/options_pb.js
echo "export {}" > $GEN_PATH/authoption/options_pb.d.ts
