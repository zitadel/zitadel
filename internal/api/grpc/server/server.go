package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	grpc_api "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

type Server interface {
	Gateway
	RegisterServer(*grpc.Server)
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

func CreateServer(verifier *authz.TokenVerifier, authConfig authz.Config, queries *query.Queries, hostHeaderName string) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.InstanceInterceptor(queries, hostHeaderName, system_pb.SystemService_MethodPrefix),
				middleware.AuthorizationInterceptor(verifier, authConfig),
				middleware.TranslationHandler(),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}
