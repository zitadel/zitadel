package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_api "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	"github.com/caos/zitadel/internal/telemetry/metrics"
)

type Server interface {
	Gateway
	RegisterServer(*grpc.Server)
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

func CreateServer(verifier *authz.TokenVerifier, authConfig authz.Config, lang language.Tag) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.AuthorizationInterceptor(verifier, authConfig),
				middleware.TranslationHandler(lang),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}
