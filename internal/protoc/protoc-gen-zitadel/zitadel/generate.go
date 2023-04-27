package zitadel

//go:generate protoc -I. -I$GOPATH/src --go-grpc_out=$GOPATH/src options.proto
