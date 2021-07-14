package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/caos/zitadel/internal/statik"
)

func ValidationHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return validate(ctx, req, info, handler)
	}
}

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
