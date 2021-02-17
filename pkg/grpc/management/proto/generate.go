package proto

//go:generate protoc -I${GOPATH}/src -I../proto -I${GOPATH}/src/github.com/caos/zitadel/pkg/grpc/message -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate -I${GOPATH}/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption --go-grpc_out=plugins=grpc:${GOPATH}/src --grpc-gateway_out=logtostderr=true,allow_delete_body=true:${GOPATH}/src --swagger_out=logtostderr=true,allow_delete_body=true:.. --authoption_out=.. --validate_out=lang=go:${GOPATH}/src management.proto
//go:generate mockgen -package api -destination ../mock/management.proto.mock.go github.com/caos/zitadel/pkg/grpc/management ManagementServiceClient
