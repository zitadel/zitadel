package middleware

import (
	"strings"

	grpc_trace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/stats"

	grpc_utils "github.com/zitadel/zitadel/internal/api/grpc"
)

type GRPCMethod string

func DefaultTracingClient() stats.Handler {
	return TracingClient(grpc_utils.Healthz, grpc_utils.Readiness, grpc_utils.Validation)
}

func TracingClient(ignoredMethods ...GRPCMethod) stats.Handler {
	return grpc_trace.NewClientHandler(grpc_trace.WithFilter(
		func(info *stats.RPCTagInfo) bool {
			for _, ignoredMethod := range ignoredMethods {
				if strings.HasSuffix(info.FullMethodName, string(ignoredMethod)) {
					return false
				}
			}
			return true
		},
	))
}
