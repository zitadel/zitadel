package proto

//go:generate protoc -I$GOPATH/src -I. --go-grpc_out=plugins=grpc:$GOPATH/src ./message.proto
