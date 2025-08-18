package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	_ "github.com/zitadel/zitadel/internal/statik"
)

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return toGRPCError(ctx, req, handler)
	}
}

func toGRPCError(ctx context.Context, req any, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	return resp, gerrors.ZITADELToGRPCError(err)
}
