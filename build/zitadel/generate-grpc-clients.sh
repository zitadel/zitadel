#! /bin/sh

if [ -n $1 ]; then
    GO_MESSAGE_IMPORT=$1/zitadel/message
else
    echo "need message import"
    exit 3
fi

protoc \
    -I=/proto/include \
    -I=/proto/include/zitadel \
    --go_opt=Mzitadel/message.proto=${GO_MESSAGE_IMPORT} \
    --go_out /go/src \
    --go-grpc_out /go/src \
    $(find /proto/include/zitadel -iname *.proto)
