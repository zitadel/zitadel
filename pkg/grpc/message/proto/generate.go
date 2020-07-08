package proto

//go:generate protoc -I$GOPATH/src -I. --go_out=plugins=grpc:$GOPATH/src ./message.proto
