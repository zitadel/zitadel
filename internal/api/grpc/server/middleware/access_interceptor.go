package middleware

import (
	"context"

	_ "github.com/zitadel/zitadel/internal/statik"
	"google.golang.org/grpc"
)

func AccessInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, err
	}
}
