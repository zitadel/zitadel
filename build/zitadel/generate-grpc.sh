#! /bin/sh

set -eux

echo "Generate grpc"

protoc \
  -I=.tmp/protos/message \
  -I=.tmp/protos/admin/proto \
  -I=.tmp/protos/management/proto \
  -I=.tmp/protos/auth/proto \
  -I=.tmp/protos \
  -I=${GOPATH}/src \
  --go_out=plugins=grpc:$GOPATH/src \
  .tmp/protos/message/proto/message.proto

protoc \
  -I=.tmp/protos/message \
  -I=.tmp/protos/admin/proto \
  -I=.tmp/protos/management/proto \
  -I=.tmp/protos/auth/proto \
  -I=.tmp/protos \
  -I=${GOPATH}/src \
  --go_out=plugins=grpc:$GOPATH/src \
  --grpc-gateway_out=logtostderr=true:$GOPATH/src \
  --swagger_out=logtostderr=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  .tmp/protos/admin/proto/admin.proto

mv admin* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/admin/

protoc \
  -I=.tmp/protos/message \
  -I=.tmp/protos/admin/proto \
  -I=.tmp/protos/management/proto \
  -I=.tmp/protos/auth/proto \
  -I=.tmp/protos \
  -I=${GOPATH}/src \
  --go_out=plugins=grpc:$GOPATH/src \
  --grpc-gateway_out=logtostderr=true,allow_delete_body=true:${GOPATH}/src \
  --swagger_out=logtostderr=true,allow_delete_body=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  .tmp/protos/management/proto/management.proto

mv management* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/management/

protoc \
  -I=.tmp/protos/message \
  -I=.tmp/protos/admin/proto \
  -I=.tmp/protos/management/proto \
  -I=.tmp/protos/auth/proto \
  -I=.tmp/protos \
  -I=${GOPATH}/src \
  --go_out=plugins=grpc:$GOPATH/src \
  --grpc-gateway_out=logtostderr=true:$GOPATH/src \
  --swagger_out=logtostderr=true:. \
  --authoption_out=. \
  --validate_out=lang=go:${GOPATH}/src \
  .tmp/protos/auth/proto/auth.proto

mv auth* $GOPATH/src/github.com/caos/zitadel/pkg/grpc/auth/

echo "done generating"
