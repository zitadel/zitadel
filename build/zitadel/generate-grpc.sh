#! /bin/sh

set -eux

echo "Generate grpc"

# output folder for openapi v2
mkdir -p $GOPATH/src/github.com/caos/zitadel/openapi/v2

ls /proto/include

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  /proto/include/zitadel/message.proto

protoc \
  -I=/proto/include \
  --go_out ${GOPATH}/src \
  --go-grpc_out ${GOPATH}/src \
  --go-grpc_opt paths=source_relative \
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/admin.proto
  # --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  # --grpc-gateway_opt logtostderr=true \
  # --openapiv2_out $GOPATH/src/github.com/caos/zitadel/openapi/v2 \
  # --openapiv2_opt logtostderr=true \ 

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  --grpc-gateway_opt logtostderr=true,paths=source_relative \
  --grpc-gateway_opt allow_delete_body=true \
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/management.proto
  # --openapiv2_out $GOPATH/src/github.com/caos/zitadel/openapi/v2 \
  # --openapiv2_opt logtostderr=true \ 

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out $GOPATH/src/github.com/caos/zitadel/openapi/v2 \
  --openapiv2_opt logtostderr=true \ 
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/auth.proto

echo "done generating grpc"