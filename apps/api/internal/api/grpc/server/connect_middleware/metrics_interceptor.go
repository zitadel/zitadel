package connect_middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

const (
	GrpcMethod                         = "grpc_method"
	ReturnCode                         = "return_code"
	GrpcRequestCounter                 = "grpc.server.request_counter"
	GrpcRequestCounterDescription      = "Grpc request counter"
	TotalGrpcRequestCounter            = "grpc.server.total_request_counter"
	TotalGrpcRequestCounterDescription = "Total grpc request counter"
	GrpcStatusCodeCounter              = "grpc.server.grpc_status_code"
	GrpcStatusCodeCounterDescription   = "Grpc status code counter"
)

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethodSuffixes ...string) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return RegisterMetrics(ctx, req, handler, metricTypes, ignoredMethodSuffixes...)
		}
	}
}

func RegisterMetrics(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, metricTypes []metrics.MetricType, ignoredMethodSuffixes ...string) (_ connect.AnyResponse, err error) {
	if len(metricTypes) == 0 {
		return handler(ctx, req)
	}

	for _, ignore := range ignoredMethodSuffixes {
		if strings.HasSuffix(req.Spec().Procedure, ignore) {
			return handler(ctx, req)
		}
	}

	resp, err := handler(ctx, req)
	if containsMetricsMethod(metrics.MetricTypeRequestCount, metricTypes) {
		RegisterGrpcRequestCounter(ctx, req.Spec().Procedure)
	}
	if containsMetricsMethod(metrics.MetricTypeTotalCount, metricTypes) {
		RegisterGrpcTotalRequestCounter(ctx)
	}
	if containsMetricsMethod(metrics.MetricTypeStatusCode, metricTypes) {
		RegisterGrpcRequestCodeCounter(ctx, req.Spec().Procedure, err)
	}
	return resp, err
}

func RegisterGrpcRequestCounter(ctx context.Context, path string) {
	var labels = map[string]attribute.Value{
		GrpcMethod: attribute.StringValue(path),
	}
	err := metrics.RegisterCounter(GrpcRequestCounter, GrpcRequestCounterDescription)
	logging.OnError(err).Warn("failed to register grpc request counter")
	err = metrics.AddCount(ctx, GrpcRequestCounter, 1, labels)
	logging.OnError(err).Warn("failed to add grpc request count")
}

func RegisterGrpcTotalRequestCounter(ctx context.Context) {
	err := metrics.RegisterCounter(TotalGrpcRequestCounter, TotalGrpcRequestCounterDescription)
	logging.OnError(err).Warn("failed to register total grpc request counter")
	err = metrics.AddCount(ctx, TotalGrpcRequestCounter, 1, nil)
	logging.OnError(err).Warn("failed to add total grpc request count")
}

func RegisterGrpcRequestCodeCounter(ctx context.Context, path string, err error) {
	statusCode := connect.CodeOf(err)
	var labels = map[string]attribute.Value{
		GrpcMethod: attribute.StringValue(path),
		ReturnCode: attribute.IntValue(runtime.HTTPStatusFromCode(codes.Code(statusCode))),
	}
	err = metrics.RegisterCounter(GrpcStatusCodeCounter, GrpcStatusCodeCounterDescription)
	logging.OnError(err).Warn("failed to register grpc status code counter")
	err = metrics.AddCount(ctx, GrpcStatusCodeCounter, 1, labels)
	logging.OnError(err).Warn("failed to add grpc status code count")
}

func containsMetricsMethod(metricType metrics.MetricType, metricTypes []metrics.MetricType) bool {
	for _, m := range metricTypes {
		if m == metricType {
			return true
		}
	}
	return false
}
