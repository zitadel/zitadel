package server

import (
	"crypto/tls"
	"net/http"

	"connectrpc.com/connect"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	grpc_api "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

type Server interface {
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

type GrpcServer interface {
	Server
	RegisterServer(*grpc.Server)
}

type WithGateway interface {
	Server
	RegisterGateway() RegisterGatewayFunc
}

// WithGatewayPrefix extends the server interface with a prefix for the grpc gateway
//
// it's used for the System, Admin, Mgmt and Auth API
type WithGatewayPrefix interface {
	GrpcServer
	WithGateway
	GatewayPathPrefix() string
}

type ConnectServer interface {
	Server
	RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler)
	FileDescriptor() protoreflect.FileDescriptor
}

func CreateServer(
	verifier authz.APITokenVerifier,
	systemAuthz authz.Config,
	authConfig authz.Config,
	queries *query.Queries,
	externalDomain string,
	tlsConfig *tls.Config,
	accessSvc *logstore.Service[*record.AccessLog],
	targetEncAlg crypto.EncryptionAlgorithm,
	translator *i18n.Translator,
) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.CallDurationHandler(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.NoCacheInterceptor(),
				middleware.InstanceInterceptor(queries, externalDomain, translator, system_pb.SystemService_ServiceDesc.ServiceName, healthpb.Health_ServiceDesc.ServiceName),
				middleware.AccessStorageInterceptor(accessSvc),
				middleware.ErrorHandler(),
				middleware.LimitsInterceptor(system_pb.SystemService_ServiceDesc.ServiceName),
				middleware.AuthorizationInterceptor(verifier, systemAuthz, authConfig),
				middleware.TranslationHandler(),
				middleware.QuotaExhaustedInterceptor(accessSvc, system_pb.SystemService_ServiceDesc.ServiceName),
				middleware.ExecutionHandler(targetEncAlg),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
				middleware.ActivityInterceptor(),
			),
		),
		grpc.StatsHandler(middleware.DefaultTracingServer()),
	}
	if tlsConfig != nil {
		serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}
	return grpc.NewServer(serverOptions...)
}
