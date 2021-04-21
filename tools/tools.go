// +build tools

package tools

import (
	//proto
	_ "github.com/envoyproxy/protoc-gen-validate"
	//proto custom authoptions
	_ "github.com/go-bindata/go-bindata/v3/go-bindata"
	// gateway grpc to rest
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	// openapi v2 descriptions
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	// grpc generator
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	//protoc
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	//generate static files
	_ "github.com/rakyll/statik"
	//proto
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	//proto doc
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	//mock
	_ "github.com/golang/mock/mockgen"
)
