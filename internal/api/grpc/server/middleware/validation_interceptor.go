package middleware

import (
	"context"

	//import to make sure go.mod does not lose it
	//because dependency is only needed for generated code
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return validate(ctx, req, info, handler)
	}
}

// validator interface needed for github.com/envoyproxy/protoc-gen-validate
// (it does not expose an interface itself)
type validator interface {
	Validate() error
}

func validate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	validate, ok := req.(validator)
	if !ok {
		return handler(ctx, req)
	}
	err := validate.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return handler(ctx, req)
}
