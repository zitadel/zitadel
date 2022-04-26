package middleware

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/errors"

	"google.golang.org/grpc"

	_ "github.com/zitadel/zitadel/internal/statik"
)

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return toGRPCError(ctx, req, handler)
	}
}

func toGRPCError(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	return resp, errors.CaosToGRPCError(ctx, err)
}
