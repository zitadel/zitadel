package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/service"
	_ "github.com/zitadel/zitadel/internal/statik"
)

func CallTimeHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = service.WithCallTimeNow(ctx)
		return handler(ctx, req)
	}
}
