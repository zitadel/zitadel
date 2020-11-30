package middleware

import (
	"context"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/caos/zitadel/internal/api/grpc/errors"
	"google.golang.org/grpc"

	_ "github.com/caos/zitadel/internal/statik"
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

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return RegisterMetrics(ctx, req, info, handler, metricTypes, ignoredMethods...)
	}
}

func RegisterMetrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, metricTypes []metrics.MetricType, ignoredMethods ...string) (_ interface{}, err error) {
	if len(metricTypes) == 0 {
		handler(ctx, req)
		return
	}

	for _, ignore := range ignoredMethods {
		if info.FullMethod == ignore {
			handler(ctx, req)
			return
		}
	}

	resp, err := handler(ctx, req)
	if containsMetricsMethod(metrics.MetricTypeRequestCount, metricTypes) {
		RegisterGrpcRequestCounter(ctx, info)
	}
	if containsMetricsMethod(metrics.MetricTypeTotalCount, metricTypes) {
		RegisterGrpcTotalRequestCounter(ctx)
	}
	if containsMetricsMethod(metrics.MetricTypeStatusCode, metricTypes) {
		RegisterGrpcRequestCodeCounter(ctx, info, err)
	}
	return resp, err
}

func RegisterGrpcRequestCounter(ctx context.Context, info *grpc.UnaryServerInfo) {
	var labels = map[string]interface{}{
		GrpcMethod: info.FullMethod,
	}
	metrics.RegisterCounter(GrpcRequestCounter, GrpcRequestCounterDescription)
	metrics.AddCount(ctx, GrpcRequestCounter, 1, labels)
}

func RegisterGrpcTotalRequestCounter(ctx context.Context) {
	metrics.RegisterCounter(TotalGrpcRequestCounter, TotalGrpcRequestCounterDescription)
	metrics.AddCount(ctx, TotalGrpcRequestCounter, 1, nil)
}

func RegisterGrpcRequestCodeCounter(ctx context.Context, info *grpc.UnaryServerInfo, err error) {
	statusCode, _, _, _ := errors.ExtractCaosError(err)
	var labels = map[string]interface{}{
		GrpcMethod: info.FullMethod,
		ReturnCode: runtime.HTTPStatusFromCode(statusCode),
	}
	metrics.RegisterCounter(GrpcStatusCodeCounter, GrpcStatusCodeCounterDescription)
	metrics.AddCount(ctx, GrpcStatusCodeCounter, 1, labels)
}

func containsMetricsMethod(metricType metrics.MetricType, metricTypes []metrics.MetricType) bool {
	for _, m := range metricTypes {
		if m == metricType {
			return true
		}
	}
	return false
}
