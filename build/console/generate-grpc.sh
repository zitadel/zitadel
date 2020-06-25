#! /bin/sh

set -eux

GEN_PATH=${GOPATH}/src/github.com/caos/zitadel/console/src/app/proto/generated

echo "Remove old files"
rm -rf $GEN_PATH

echo "Create folders"
mkdir -p $GEN_PATH

echo "Generate grpc"


mkdir -p ./validate/validate
touch ./validate/validate.proto

mkdir -p ./grpc-ecosystem/grpc-gateway
touch ./grpc-ecosystem/grpc-gateway/validate.proto

mkdir -p ./protoc-gen-swagger/options
touch ./protoc-gen-swagger/options/annotations.proto

protoc \
  -I=/usr/local/include \
  -I=${GOPATH}/src/github.com/caos/zitadel/pkg/management/grpc/proto \
  -I=${GOPATH}/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption \
  -I=${GOPATH}/src/github.com/caos/zitadel/console/node_modules/google-proto-files \
  -I=. \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  ${GOPATH}/src/github.com/caos/zitadel/pkg/management/grpc/proto/*.proto \


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
