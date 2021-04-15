#! /bin/sh

set -eux

echo "Generate grpc"

OPENAPI_PATH=${GOPATH}/src/github.com/caos/zitadel/openapi/v2
ZITADEL_PATH=${GOPATH}/src/github.com/caos/zitadel
GRPC_PATH=${ZITADEL_PATH}/pkg/grpc
PROTO_PATH=/proto/include/zitadel
DOCS_PATH=${ZITADEL_PATH}/docs/apis/proto

# generate go stub and grpc code for all files
protoc \
  -I=/proto/include/ \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  $(find ${PROTO_PATH} -iname *.proto)

# generate authoptions code from templates
go-bindata \
  -pkg main \
  -prefix internal/protoc/protoc-gen-authoption \
  -o ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption/templates.gen.go \
  ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption/templates

# install authoption proto compiler
go install ${ZITADEL_PATH}/internal/protoc/protoc-gen-authoption

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
  --authoption_out ${GRPC_PATH}/admin \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/admin.proto

# authoptions are generated into the wrong folder
mv ${ZITADEL_PATH}/pkg/grpc/admin/zitadel/* ${ZITADEL_PATH}/pkg/grpc/admin
rm -r ${ZITADEL_PATH}/pkg/grpc/admin/zitadel

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt allow_delete_body=true \
  --authoption_out ${GRPC_PATH}/management \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/management.proto

# authoptions are generated into the wrong folder
mv ${ZITADEL_PATH}/pkg/grpc/management/zitadel/* ${ZITADEL_PATH}/pkg/grpc/management
rm -r ${ZITADEL_PATH}/pkg/grpc/management/zitadel

protoc \
  -I=/proto/include \
  --grpc-gateway_out ${GOPATH}/src \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out ${OPENAPI_PATH} \
  --openapiv2_opt logtostderr=true \
  --authoption_out=${GRPC_PATH}/auth \
  --validate_out=lang=go:${GOPATH}/src \
  ${PROTO_PATH}/auth.proto

# authoptions are generated into the wrong folder
mv ${ZITADEL_PATH}/pkg/grpc/auth/zitadel/* ${ZITADEL_PATH}/pkg/grpc/auth
rm -r ${ZITADEL_PATH}/pkg/grpc/auth/zitadel

## generate docs
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,auth.md \
  ${PROTO_PATH}/auth.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,management.md \
  ${PROTO_PATH}/management.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,admin.md \
  ${PROTO_PATH}/admin.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,app.md \
  ${PROTO_PATH}/app.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,auth_n_key.md \
  ${PROTO_PATH}/auth_n_key.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,change.md \
  ${PROTO_PATH}/change.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,idp.md \
  ${PROTO_PATH}/idp.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,member.md \
  ${PROTO_PATH}/member.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,message.md \
  ${PROTO_PATH}/message.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,object.md \
  ${PROTO_PATH}/object.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,options.md \
  ${PROTO_PATH}/options.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,org.md \
  ${PROTO_PATH}/org.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,policy.md \
  ${PROTO_PATH}/policy.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,project.md \
  ${PROTO_PATH}/project.proto
protoc \
  -I=/proto/include \
  --doc_out=${DOCS_PATH} --doc_opt=${PROTO_PATH}/docs/zitadel-md.tmpl,user.md \
  ${PROTO_PATH}/user.proto

echo "done generating grpc"