package middleware

import (
	"context"
	"strings"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	grpc_utils "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/tracing"
)

type GRPCMethod string

func TracingStatsServer(ignoredMethods ...GRPCMethod) grpc.ServerOption {
	return grpc.StatsHandler(
		&tracingServerHandler{
			ignoredMethods,
			ocgrpc.ServerHandler{
				StartOptions: trace.StartOptions{
					Sampler:  tracing.Sampler(),
					SpanKind: trace.SpanKindServer,
				},
			},
		},
	)
}

func DefaultTracingStatsServer() grpc.ServerOption {
	return TracingStatsServer(grpc_utils.Healthz, grpc_utils.Readiness, grpc_utils.Validation)
}

type tracingServerHandler struct {
	IgnoredMethods []GRPCMethod
	ocgrpc.ServerHandler
}

func (s *tracingServerHandler) TagRPC(ctx context.Context, tagInfo *stats.RPCTagInfo) context.Context {
	for _, method := range s.IgnoredMethods {
		if strings.HasSuffix(tagInfo.FullMethodName, string(method)) {
			return ctx
		}
	}
	return s.ServerHandler.TagRPC(ctx, tagInfo)
}
