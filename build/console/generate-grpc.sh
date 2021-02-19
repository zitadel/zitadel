#! /bin/sh

set -eux

GEN_PATH=src/app/proto/generated

echo "Create folders"
mkdir -p $GEN_PATH

echo "Generate grpc"

protoc \
  -I=/proto/include \
  -I=node_modules/google-proto-files \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  $(find /proto/include -iname "*.proto")