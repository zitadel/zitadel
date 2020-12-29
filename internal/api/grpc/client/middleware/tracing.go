package middleware

import (
	"context"
	"strings"

	grpc_utils "github.com/caos/zitadel/internal/api/grpc"
	grpc_trace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type GRPCMethod string

func DefaultTracingClient() grpc.UnaryClientInterceptor {
	return TracingServer(grpc_utils.Healthz, grpc_utils.Readiness, grpc_utils.Validation)
}

func TracingServer(ignoredMethods ...GRPCMethod) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		for _, ignoredMethod := range ignoredMethods {
			if strings.HasSuffix(method, string(ignoredMethod)) {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
		}
		return grpc_trace.UnaryClientInterceptor()(ctx, method, req, reply, cc, invoker, opts...)
	}
}
