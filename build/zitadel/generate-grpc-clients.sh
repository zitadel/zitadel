#! /bin/sh

if [ -n $1 ]; then
    GO_MESSAGE_IMPORT=$1/zitadel/message
else
    echo "need message import"
    exit 3
fi

generate () {
    protoc \
        -I=/.tmp/protos \
        -I=/go/src/github.com/caos/zitadel/pkg/grpc/message \
        -I=/go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption \
        -I=/go/src \
        --go_opt=Mproto/message.proto=${GO_MESSAGE_IMPORT} \
        --go_out /go/src \
        --go-grpc_out /go/src \
        $1/$2
}
generate /go/src/github.com/caos/zitadel/pkg/grpc/message/proto message.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/admin/proto admin.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/auth/proto auth.proto
generate /go/src/github.com/caos/zitadel/pkg/grpc/management/proto management.proto
