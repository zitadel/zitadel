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
  $(find ${PROTO_PATH} -iname *.proto | grep -v "management|admin|auth")

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
  --go_out ${GOPATH}/src \
  --go-grpc_out ${GOPATH}/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out ${GRPC_PATH}/admin \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/admin.proto
mv ${ZITADEL_PATH}/pkg/grpc/admin/zitadel/* ${ZITADEL_PATH}/pkg/grpc/admin
rm -r ${ZITADEL_PATH}/pkg/grpc/admin/zitadel

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --authoption_out ${GRPC_PATH}/management \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/management.proto
mv ${ZITADEL_PATH}/pkg/grpc/management/zitadel/* ${ZITADEL_PATH}/pkg/grpc/management
rm -r ${ZITADEL_PATH}/pkg/grpc/management/zitadel

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GRPC_PATH}/auth \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/auth.proto
mv ${ZITADEL_PATH}/pkg/grpc/auth/zitadel/* ${ZITADEL_PATH}/pkg/grpc/auth
rm -r ${ZITADEL_PATH}/pkg/grpc/auth/zitadel

echo "done generating grpc"