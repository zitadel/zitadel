// +build tools

package tools

import (
	//proto
	_ "github.com/envoyproxy/protoc-gen-validate"
	//proto custom authoptions
	_ "github.com/go-bindata/go-bindata/go-bindata"
	//proto
	_ "github.com/golang/protobuf/protoc-gen-go"
	//proto gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	//proto gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	//generate static files
	_ "github.com/rakyll/statik"
	//proto
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
