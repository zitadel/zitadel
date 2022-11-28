package middleware

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/service"
	_ "github.com/zitadel/zitadel/internal/statik"
	"google.golang.org/grpc"
)

func CallTimeHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = service.WithCallTimeNow(ctx)
		return handler(ctx, req)
	}
}
