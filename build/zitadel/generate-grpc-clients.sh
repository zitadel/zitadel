#! /bin/sh

if [ -n $1 ]; then
    GO_MESSAGE_IMPORT=$1/zitadel/message
else
    echo "need message import"
    exit 3
fi

PROTO_PATH=/go/src/github.com/caos/zitadel/proto

protoc \
    -I=/.tmp/protos \
    -I=/go/src/github.com/caos/zitadel/pkg/grpc/message \
    -I=/go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption \
    -I=/go/src \
    --go_opt=Mproto/message.proto=${GO_MESSAGE_IMPORT} \
    --go_out /go/src \
    --go-grpc_out /go/src \
    $(find ${PROTO_PATH} -iname *.proto)
