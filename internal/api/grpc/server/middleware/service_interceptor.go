package middleware

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/service"
	_ "github.com/zitadel/zitadel/internal/statik"
	"google.golang.org/grpc"
)

func ServiceHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		namer := info.Server.(interface{ AppName() string })
		ctx = service.WithService(ctx, namer.AppName())
		return handler(ctx, req)
	}
}
