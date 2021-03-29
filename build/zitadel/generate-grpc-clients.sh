#! /bin/sh

if [ -n $1 ]; then
    GO_OUT=$1
else
    GO_OUT=$GOPATH/src
fi

if [ -n $2 ]; then
    GRPC_OUT=$2
else
    GRPC_OUT=$GO_OUT
fi

if [ -n $3 ]; then
    GO_MESSAGE_IMPORT=$3
else
    echo "need message import"
    exit 3
fi

generate () {
    protoc \
        -I=/.tmp/protos \
        -I=/go/src/github.com/caos/zitadel/pkg/grpc/message \
        -I=/go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption \
        -I=$1 \
        --go_opt=Mproto/message.proto=${GO_MESSAGE_IMPORT} \
        --go_out ${GO_OUT} \
        --go-grpc_out ${GO_OUT} \
        $1/$2
}
generate /go/src/github.com/caos/zitadel/pkg/grpc/message/proto message.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/admin/proto admin.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/auth/proto auth.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/management/proto management.proto
