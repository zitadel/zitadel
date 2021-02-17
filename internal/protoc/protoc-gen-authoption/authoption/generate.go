package authoption

//go:generate protoc -I. -I$GOPATH/src --go-grpc_out=plugins=grpc:$GOPATH/src options.proto
