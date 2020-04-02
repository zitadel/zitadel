package middleware

import (
	"context"
	"strings"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"

	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/tracing"
)

type GRPCMethod string

func TracingStatsClient(ignoredMethods ...GRPCMethod) grpc.DialOption {
	return grpc.WithStatsHandler(&tracingClientHandler{ignoredMethods, ocgrpc.ClientHandler{StartOptions: trace.StartOptions{Sampler: tracing.Sampler(), SpanKind: trace.SpanKindClient}}})
}

func DefaultTracingStatsClient() grpc.DialOption {
	return TracingStatsClient(api.Healthz, api.Readiness, api.Validation)
}

type tracingClientHandler struct {
	IgnoredMethods []GRPCMethod
	ocgrpc.ClientHandler
}

func (s *tracingClientHandler) TagRPC(ctx context.Context, tagInfo *stats.RPCTagInfo) context.Context {
	for _, method := range s.IgnoredMethods {
		if strings.HasSuffix(tagInfo.FullMethodName, string(method)) {
			return ctx
		}
	}
	return s.ClientHandler.TagRPC(ctx, tagInfo)
}
