#! /bin/sh

if [ -n $1 ]; then
    GO_MESSAGE_IMPORT=$1/zitadel/message
else
    echo "need message import"
    exit 3
fi

protoc \
    -I=/protos/include \
    --go_opt=Mproto/message.proto=${GO_MESSAGE_IMPORT} \
    --go_out /go/src \
    --go-grpc_out /go/src \
    $(find /proto/include/zitadel -iname *.proto)
