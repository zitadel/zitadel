package api

import (
	"context"
	"net/http"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_api "github.com/caos/zitadel/internal/api/grpc"
	admin_grpc "github.com/caos/zitadel/internal/api/grpc/admin"
	auth_grpc "github.com/caos/zitadel/internal/api/grpc/auth"
	mgmt_grpc "github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/caos/zitadel/v2/api/admin"
	"github.com/caos/zitadel/v2/api/auth"
	"github.com/caos/zitadel/v2/api/mgmt"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
)

var (
	gwOpts = server.DefaultServeMuxOptions("x-zitadel-")
)

type API struct{}

func New(ctx context.Context, baseRouter *mux.Router, mgmtSrv *mgmt_grpc.Server, adminSrv *admin_grpc.Server, authSrv *auth_grpc.Server, verifier *authz.TokenVerifier, authZConf authz.Config) *API {
	grpcHandler := grpcServer(verifier, authZConf)
	grpcWebHandler := grpcweb.WrapServer(grpcHandler)
	apiRoute := baseRouter.PathPrefix("/api").Subrouter()

	// services
	mgmtHandler := mgmt.New(mgmtSrv)
	routeAPI(ctx, mgmtHandler, grpcHandler, apiRoute)
	adminHandler := admin.New(adminSrv)
	routeAPI(ctx, adminHandler, grpcHandler, apiRoute)
	authHandler := auth.New(authSrv)
	routeAPI(ctx, authHandler, grpcHandler, apiRoute)

	routeGRPC(baseRouter, grpcHandler, grpcWebHandler)

	return &API{}
}

type handler interface {
	RegisterRESTGateway(context.Context, *runtime.ServeMux) error
	RegisterGRPC(*grpc.Server)
	ServicePrefix() string
}

func routeAPI(ctx context.Context, h handler, grpcServ *grpc.Server, apiRouter *mux.Router) error {
	m := runtime.NewServeMux(gwOpts...)
	h.RegisterGRPC(grpcServ)
	if err := h.RegisterRESTGateway(ctx, m); err != nil {
		panic(err)
	}
	apiRouter.PathPrefix(h.ServicePrefix()).Handler(http.StripPrefix("/api"+h.ServicePrefix(), m))
	return nil
}

func routeGRPC(baseRouter *mux.Router, grpcHandler *grpc.Server, grpcWebHandler *grpcweb.WrappedGrpcServer) {
	http2Route := baseRouter.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
		return r.ProtoMajor == 2
	}).Subrouter()
	http2Route.Headers("Content-Type", "application/grpc").Handler(grpcHandler)
	http2Route.NewRoute().HeadersRegexp("Content-Type", "application/grpc-web.*").Handler(grpcWebHandler)
}

func grpcServer(verifier *authz.TokenVerifier, authZConf authz.Config) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(

		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.AuthorizationInterceptor(verifier, authZConf),
				middleware.TranslationHandler(language.German),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}
