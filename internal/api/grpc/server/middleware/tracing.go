package middleware

import (
	"context"
	"strings"

	grpc_utils "github.com/caos/zitadel/internal/api/grpc"
	grpc_trace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type GRPCMethod string

func DefaultTracingServer() grpc.UnaryServerInterceptor {
	return TracingServer(grpc_utils.Healthz, grpc_utils.Readiness, grpc_utils.Validation)
}

func TracingServer(ignoredMethods ...GRPCMethod) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		for _, ignoredMethod := range ignoredMethods {
			if strings.HasSuffix(info.FullMethod, string(ignoredMethod)) {
				return handler(ctx, req)
			}
		}
		return grpc_trace.UnaryServerInterceptor()(ctx, req, info, handler)
	}
}
