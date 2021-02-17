#! /bin/sh

set -eux

echo "Generate grpc"

OPENAPI_PATH=${GOPATH}/src/github.com/caos/zitadel/openapi/v2
GRPC_PATH=${GOPATH}/src/github.com/caos/zitadel/pkg/grpc
PROTO_PATH=/proto/include/zitadel

# output folder for openapi v2
mkdir -p ${OPENAPI_PATH}

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  ${PROTO_PATH}/message.proto

protoc \
  -I=/proto/include \
  --go_out ${GOPATH}/src \
  --go-grpc_out ${GOPATH}/src \
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ${GRPC_PATH} \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GRPC_PATH} \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/admin.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ${GRPC_PATH} \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GRPC_PATH} \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/management.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ${GRPC_PATH} \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GRPC_PATH} \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/auth.proto

echo "done generating grpc"