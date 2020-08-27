package middleware

import (
	"context"

	"google.golang.org/grpc"

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
		return nil, err
	}
	return handler(ctx, req)
}
