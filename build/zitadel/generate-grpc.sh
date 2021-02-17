#! /bin/sh

set -eux

echo "Generate grpc"

OPENAPI_PATH=${GOPATH}/src/github.com/caos/zitadel/openapi/v2
ZITADEL_PATH=${GOPATH}/src/github.com/caos/zitadel
GRPC_PATH=${ZITADEL_PATH}/pkg/grpc
PROTO_PATH=/proto/include/zitadel

protoc \
  -I=/proto/include/ \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  ${PROTO_PATH}/options.proto

go-bindata \
  -pkg main \
  -prefix internal/protoc/protoc-gen-authoption \
  -o ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption/templates.gen.go \
  ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption/templates

go install ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption

# output folder for openapi v2
mkdir -p ${OPENAPI_PATH}

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  ${PROTO_PATH}/message.proto

protoc \
  -I=/proto/include \
  --go_out ${GOPATH}/src \
  --go-grpc_out ${GOPATH}/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/admin.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/management.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GOPATH}/src \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/auth.proto

echo "done generating grpc"