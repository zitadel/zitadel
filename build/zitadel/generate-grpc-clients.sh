#! /bin/sh

echo "Generate grpc"

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

generate /go/src/github.com/caos/zitadel/pkg/message/proto message.proto
# generate /go/src/github.com/caos/zitadel/pkg/admin/proto admin.proto
# generate /go/src/github.com/caos/zitadel/pkg/auth/proto auth.proto
# generate /go/src/github.com/caos/zitadel/pkg/management/proto management.proto

echo "done generating grpc"