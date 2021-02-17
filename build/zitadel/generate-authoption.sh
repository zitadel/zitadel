#! /bin/sh

set -eux

echo "Generate authoption"

protoc \
    -I=/proto/include/ \
    --go_out $GOPATH/src \
    --go-grpc_out $GOPATH/src \
    /proto/include/zitadel/options.proto

echo "done generate authoption" 