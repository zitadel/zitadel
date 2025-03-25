package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/service"
	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func ServiceHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, span := tracing.NewSpan(ctx)
		defer span.End()

		namer, ok := info.Server.(interface{ AppName() string })
		if !ok {
			return handler(ctx, req)
		}
		ctx = service.WithService(ctx, namer.AppName())
		return handler(ctx, req)
	}
}
