#!/bin/sh

set -eux

echo "Generate grpc"

OPENAPI_PATH=${GOPATH}/src/github.com/zitadel/zitadel/openapi/v2
ZITADEL_PATH=${GOPATH}/src/github.com/zitadel/zitadel
GRPC_PATH=${ZITADEL_PATH}/pkg/grpc
PROTO_PATH=/proto/include/zitadel
DOCS_PATH=${ZITADEL_PATH}/docs/apis/proto

# generate go stub and grpc code for all files
protoc \
  -I=/proto/include/ \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --validate_out=lang=go:${GOPATH}/src \
  $(find ${PROTO_PATH} -iname *.proto)

# install authoption and zitadel proto compiler
go install ${ZITADEL_PATH}/internal/protoc/protoc-gen-auth
go install ${ZITADEL_PATH}/internal/protoc/protoc-gen-zitadel

# output folder for openapi v2
mkdir -p ${OPENAPI_PATH}
mkdir -p ${DOCS_PATH}

# generate additional output

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --auth_out ${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/system.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --auth_out ${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/admin.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --auth_out ${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/management.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --auth_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/auth.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --zitadel_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/user/v2alpha/user_service.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --zitadel_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/policy/v2alpha/policy_service.proto

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --auth_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/session/v2alpha/session_service.proto

echo "done generating grpc"
