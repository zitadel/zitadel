package proto

//go:generate protoc -I$GOPATH/src -I../proto -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/envoyproxy/protoc-gen-validate -I$GOPATH/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption --go-grpc_out=plugins=grpc:$GOPATH/src --grpc-gateway_out=logtostderr=true:$GOPATH/src --swagger_out=logtostderr=true:.. --authoption_out=.. --validate_out=lang=go:${GOPATH}/src admin.proto
//go:generate mockgen -package api -destination ../mock/admin.proto.mock.go github.com/caos/zitadel/pkg/grpc/admin AdminServiceClient
