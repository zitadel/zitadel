package api

import (
	"context"
	"crypto/tls"
	"net/http"
	"sort"
	"strings"

	"connectrpc.com/grpcreflect"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/zitadel/zitadel/internal/api/authz"
	grpc_api "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/api/grpc/server/connect_middleware"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

var (
	metricTypes = []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
)

type API struct {
	port              uint16
	externalDomain    string
	grpcServer        *grpc.Server
	verifier          authz.APITokenVerifier
	health            healthCheck
	router            *mux.Router
	hostHeaders       []string
	grpcGateway       *server.Gateway
	healthServer      *health.Server
	accessInterceptor *http_mw.AccessInterceptor
	queries           *query.Queries
	authConfig        authz.Config
	systemAuthZ       authz.Config
	connectServices   map[string][]string

	targetEncryptionAlgorithm crypto.EncryptionAlgorithm
	translator                *i18n.Translator
}

func (a *API) ListGrpcServices() []string {
	serviceInfo := a.grpcServer.GetServiceInfo()
	services := make([]string, len(serviceInfo)+len(a.connectServices))
	i := 0
	for servicename := range serviceInfo {
		services[i] = servicename
		i++
	}
	for prefix := range a.connectServices {
		services[i] = strings.Trim(prefix, "/")
		i++
	}
	sort.Strings(services)
	return services
}

func (a *API) ListGrpcMethods() []string {
	serviceInfo := a.grpcServer.GetServiceInfo()
	methods := make([]string, 0)
	for servicename, service := range serviceInfo {
		for _, method := range service.Methods {
			methods = append(methods, "/"+servicename+"/"+method.Name)
		}
	}
	for service, methodList := range a.connectServices {
		for _, method := range methodList {
			methods = append(methods, service+method)
		}
	}
	sort.Strings(methods)
	return methods
}

type healthCheck interface {
	Health(ctx context.Context) error
}

func New(
	ctx context.Context,
	port uint16,
	router *mux.Router,
	queries *query.Queries,
	verifier authz.APITokenVerifier,
	systemAuthz authz.Config,
	authZ authz.Config,
	tlsConfig *tls.Config,
	externalDomain string,
	hostHeaders []string,
	accessInterceptor *http_mw.AccessInterceptor,
	targetEncryptionAlgorithm crypto.EncryptionAlgorithm,
	translator *i18n.Translator,
) (_ *API, err error) {
	api := &API{
		port:                      port,
		externalDomain:            externalDomain,
		verifier:                  verifier,
		health:                    queries,
		router:                    router,
		queries:                   queries,
		accessInterceptor:         accessInterceptor,
		hostHeaders:               hostHeaders,
		authConfig:                authZ,
		systemAuthZ:               systemAuthz,
		connectServices:           make(map[string][]string),
		targetEncryptionAlgorithm: targetEncryptionAlgorithm,
		translator:                translator,
	}

	api.grpcServer = server.CreateServer(api.verifier, systemAuthz, authZ, queries, externalDomain, tlsConfig, accessInterceptor.AccessService(), targetEncryptionAlgorithm, api.translator)
	api.grpcGateway, err = server.CreateGateway(ctx, port, hostHeaders, accessInterceptor, tlsConfig)
	if err != nil {
		return nil, err
	}
	api.registerHealthServer()

	api.RegisterHandlerOnPrefix("/debug", api.healthHandler())
	api.router.Handle("/", http.RedirectHandler(login.HandlerPrefix, http.StatusFound))

	return api, nil
}

func (a *API) serverReflection() {
	reflector := grpcreflect.NewStaticReflector(a.ListGrpcServices()...)
	a.RegisterHandlerOnPrefix(grpcreflect.NewHandlerV1(reflector))
	a.RegisterHandlerOnPrefix(grpcreflect.NewHandlerV1Alpha(reflector))
}

// RegisterServer registers a grpc service on the grpc server,
// creates a new grpc gateway and registers it as a separate http handler
//
// used for v1 api (system, admin, mgmt, auth)
func (a *API) RegisterServer(ctx context.Context, grpcServer server.WithGatewayPrefix, tlsConfig *tls.Config) error {
	grpcServer.RegisterServer(a.grpcServer)
	handler, prefix, err := server.CreateGatewayWithPrefix(
		ctx,
		grpcServer,
		a.port,
		a.hostHeaders,
		a.accessInterceptor,
		tlsConfig,
	)
	if err != nil {
		return err
	}
	a.RegisterHandlerOnPrefix(prefix, handler)
	a.verifier.RegisterServer(grpcServer.AppName(), grpcServer.MethodPrefix(), grpcServer.AuthMethods())
	a.healthServer.SetServingStatus(grpcServer.MethodPrefix(), healthpb.HealthCheckResponse_SERVING)
	return nil
}

// RegisterService registers a grpc service on the grpc server,
// and its gateway on the gateway handler
//
// used for >= v2 api (e.g. user, session, ...)
func (a *API) RegisterService(ctx context.Context, srv server.Server) error {
	switch service := srv.(type) {
	case server.GrpcServer:
		service.RegisterServer(a.grpcServer)
	case server.ConnectServer:
		a.registerConnectServer(service)
	}
	if withGateway, ok := srv.(server.WithGateway); ok {
		err := server.RegisterGateway(ctx, a.grpcGateway, withGateway)
		if err != nil {
			return err
		}
	}
	a.verifier.RegisterServer(srv.AppName(), srv.MethodPrefix(), srv.AuthMethods())
	a.healthServer.SetServingStatus(srv.MethodPrefix(), healthpb.HealthCheckResponse_SERVING)
	return nil
}

func (a *API) registerConnectServer(service server.ConnectServer) {
	prefix, handler := service.RegisterConnectServer(
		connect_middleware.CallDurationHandler(),
		connect_middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
		connect_middleware.NoCacheInterceptor(),
		connect_middleware.InstanceInterceptor(a.queries, a.externalDomain, a.translator, system_pb.SystemService_ServiceDesc.ServiceName, healthpb.Health_ServiceDesc.ServiceName),
		connect_middleware.AccessStorageInterceptor(a.accessInterceptor.AccessService()),
		connect_middleware.ErrorHandler(),
		connect_middleware.LimitsInterceptor(system_pb.SystemService_ServiceDesc.ServiceName),
		connect_middleware.AuthorizationInterceptor(a.verifier, a.systemAuthZ, a.authConfig),
		connect_middleware.TranslationHandler(),
		connect_middleware.QuotaExhaustedInterceptor(a.accessInterceptor.AccessService(), system_pb.SystemService_ServiceDesc.ServiceName),
		connect_middleware.ExecutionHandler(a.targetEncryptionAlgorithm),
		connect_middleware.ValidationHandler(),
		connect_middleware.ServiceHandler(),
		connect_middleware.ActivityInterceptor(),
	)
	methods := service.FileDescriptor().Services().Get(0).Methods()
	methodNames := make([]string, methods.Len())
	for i := 0; i < methods.Len(); i++ {
		methodNames[i] = string(methods.Get(i).Name())
	}
	a.connectServices[prefix] = methodNames
	a.RegisterHandlerPrefixes(http_mw.CORSInterceptor(handler), prefix)
}

// HandleFunc allows registering a [http.HandlerFunc] on an exact
// path, instead of prefix like RegisterHandlerOnPrefix.
func (a *API) HandleFunc(path string, f http.HandlerFunc) {
	a.router.HandleFunc(path, f)
}

// RegisterHandlerOnPrefix registers a http handler on a path prefix
// the prefix will not be passed to the actual handler
func (a *API) RegisterHandlerOnPrefix(prefix string, handler http.Handler) {
	prefix = strings.TrimSuffix(prefix, "/")
	subRouter := a.router.PathPrefix(prefix).Name(prefix).Subrouter()
	subRouter.PathPrefix("").Handler(http.StripPrefix(prefix, handler))
}

// RegisterHandlerPrefixes registers a http handler on a multiple path prefixes
// the prefix will remain when calling the actual handler
func (a *API) RegisterHandlerPrefixes(handler http.Handler, prefixes ...string) {
	for _, prefix := range prefixes {
		prefix = strings.TrimSuffix(prefix, "/")
		subRouter := a.router.PathPrefix(prefix).Name(prefix).Subrouter()
		subRouter.PathPrefix("").Handler(handler)
	}
}

func (a *API) registerHealthServer() {
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(a.grpcServer, healthServer)
	a.healthServer = healthServer
}

func (a *API) RouteGRPC() {
	// since all services are now registered, we can build the grpc server reflection and register the handler
	a.serverReflection()

	http2Route := a.router.
		MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
			return r.ProtoMajor == 2
		}).
		Subrouter().
		Name("grpc")
	http2Route.
		Methods(http.MethodPost).
		HeadersRegexp(http_util.ContentType, `^application\/grpc(\+proto|\+json)?$`).
		Handler(a.grpcServer)

	a.routeGRPCWeb()
	a.router.NewRoute().
		Handler(a.grpcGateway.Handler()).
		Name("grpc-gateway")
}

func (a *API) routeGRPCWeb() {
	grpcWebServer := grpcweb.WrapServer(a.grpcServer,
		grpcweb.WithAllowedRequestHeaders(
			[]string{
				http_util.Origin,
				http_util.ContentType,
				http_util.Accept,
				http_util.AcceptLanguage,
				http_util.Authorization,
				http_util.ZitadelOrgID,
				http_util.XUserAgent,
				http_util.XGrpcWeb,
			},
		),
		grpcweb.WithOriginFunc(func(_ string) bool {
			return true
		}),
	)
	a.router.Use(http_mw.RobotsTagHandler)
	a.router.NewRoute().
		Methods(http.MethodPost, http.MethodOptions).
		MatcherFunc(
			func(r *http.Request, _ *mux.RouteMatch) bool {
				return grpcWebServer.IsGrpcWebRequest(r) || grpcWebServer.IsAcceptableGrpcCorsRequest(r)
			}).
		Handler(grpcWebServer).
		Name("grpc-web")
}

func (a *API) healthHandler() http.Handler {
	checks := []ValidationFunction{
		func(ctx context.Context) error {
			if err := a.health.Health(ctx); err != nil {
				return zerrors.ThrowInternal(err, "API-F24h2", "DB CONNECTION ERROR")
			}
			return nil
		},
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", handleHealth)
	handler.HandleFunc("/ready", handleReadiness(checks))
	handler.HandleFunc("/validate", handleValidate(checks))
	handler.Handle("/metrics", metricsExporter())

	return handler
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	logging.WithFields("traceID", tracing.TraceIDFromCtx(r.Context())).OnError(err).Error("error writing ok for health")
}

func handleReadiness(checks []ValidationFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errs := validate(r.Context(), checks)
		if len(errs) == 0 {
			http_util.MarshalJSON(w, "ok", nil, http.StatusOK)
			return
		}
		http_util.MarshalJSON(w, nil, errs[0], http.StatusPreconditionFailed)
	}
}

func handleValidate(checks []ValidationFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errs := validate(r.Context(), checks)
		if len(errs) == 0 {
			http_util.MarshalJSON(w, "ok", nil, http.StatusOK)
			return
		}
		http_util.MarshalJSON(w, errs, nil, http.StatusOK)
	}
}

type ValidationFunction func(ctx context.Context) error

func validate(ctx context.Context, validations []ValidationFunction) []error {
	errs := make([]error, 0)
	for _, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.WithFields("traceID", tracing.TraceIDFromCtx(ctx)).WithError(err).Error("validation failed")
			errs = append(errs, err)
		}
	}
	return errs
}

func metricsExporter() http.Handler {
	exporter := metrics.GetExporter()
	if exporter == nil {
		return http.NotFoundHandler()
	}
	return exporter
}
