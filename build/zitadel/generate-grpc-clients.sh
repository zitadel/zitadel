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

generate () {
    protoc \
        -I=/.tmp/protos/ \
        -I=$1 \
        --go_out ${GO_OUT} \
        --go-grpc_out $GOPATH/src \
        $1/$2
}
generate /go/src/github.com/caos/zitadel/pkg/grpc/message/proto message.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/admin/proto admin.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/auth/proto auth.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/management/proto management.proto
