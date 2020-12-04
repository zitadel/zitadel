#! /bin/sh

set -eux

GEN_PATH=src/app/proto/generated

echo "Remove old files"
rm -rf $GEN_PATH

echo "Create folders"
mkdir -p $GEN_PATH

targetcurl () {
   mkdir -p $1 && cd $1 && { curl -O $2; cd -; }
}

echo "Download additional protofiles"
targetcurl tmp/validate https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.4.0/validate/validate.proto
targetcurl tmp/protoc-gen-swagger/options https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/annotations.proto
targetcurl tmp/protoc-gen-swagger/options https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/openapiv2.proto

echo "Generate grpc"

protoc \
  -I=/usr/local/include \
  -I=../pkg/grpc/message \
  -I=../pkg/grpc/management/proto \
  -I=../pkg/grpc/auth/proto \
  -I=../pkg/grpc/admin/proto \
  -I=../internal/protoc/protoc-gen-authoption \
  -I=node_modules/google-proto-files \
  -I=tmp \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  ../pkg/grpc/message/proto/*.proto \
  ../pkg/grpc/management/proto/*.proto \
  ../pkg/grpc/admin/proto/*.proto \
  ../pkg/grpc/auth/proto/*.proto

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