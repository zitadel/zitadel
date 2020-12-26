#! /bin/sh

set -eux

echo "Generate grpc"

protoc -I=/proto/include/ -I$GOPATH/src --go_out=plugins=grpc:$GOPATH/src options.proto

protoc \
  -I=/proto/proto \
  -I=/proto/include \
  #-I=${GOPATH}/src \
  --go_grpc_out=$GOPATH/src \
  /proto/include/proto/message.proto

ls ${GOPATH}/src/github.com/caos/zitadel/pkg/grpc/message

protoc \
  -I=/proto/proto \
  -I=/proto/zitadel \
  -I=/proto/include \
  #-I=${GOPATH}/src \
  --go_grpc_out=$GOPATH/src \
  --grpc-gateway_out=logtostderr=true:$GOPATH/src \
  --swagger_out=logtostderr=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/zitadel/admin.proto

mv admin* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/admin/

protoc \
  -I=/proto/proto \
  -I=/proto/zitadel \
  -I=/proto/include \
  --go_grpc_out=$GOPATH/src \
  --grpc-gateway_out=logtostderr=true,allow_delete_body=true:${GOPATH}/src \
  --swagger_out=logtostderr=true,allow_delete_body=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/zitadel/management.proto

mv management* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/management/

protoc \
  -I=/proto/proto \
  -I=/proto/zitadel \
  -I=/proto/include \
  --go_grpc_out=$GOPATH/src \
  --grpc-gateway_out=logtostderr=true:$GOPATH/src \
  --swagger_out=logtostderr=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/zitadel/auth.proto

mv auth* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/auth/

echo "done generating"
