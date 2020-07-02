package message

//go:generate protoc -I$GOPATH/src -I. --go_out=plugins=grpc:$GOPATH/src --grpc-gateway_out=logtostderr=true,allow_delete_body=true:$GOPATH/src ./message.proto
