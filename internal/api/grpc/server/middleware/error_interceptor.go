package middleware

import (
	"context"

	"google.golang.org/grpc"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
)

func ErrorHandler() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, grpc_util.CaosToGRPCError(err)
	}
}
