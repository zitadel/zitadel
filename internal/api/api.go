package api

import (
	"context"
	"net/http"
	"strings"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"

	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type API struct {
	port           uint16
	grpcServer     *grpc.Server
	verifier       *internal_authz.TokenVerifier
	health         health
	router         *mux.Router
	externalSecure bool
}

type health interface {
	Health(ctx context.Context) error
	Instance(ctx context.Context) (*query.Instance, error)
}

func New(
	port uint16,
	router *mux.Router,
	queries *query.Queries,
	verifier *internal_authz.TokenVerifier,
	authZ internal_authz.Config,
	externalSecure bool,
	http2HostName string,
) *API {
	api := &API{
		port:           port,
		verifier:       verifier,
		health:         queries,
		router:         router,
		externalSecure: externalSecure,
	}
	api.grpcServer = server.CreateServer(api.verifier, authZ, queries, http2HostName)
	api.routeGRPC()

	api.RegisterHandler("/debug", api.healthHandler())

	return api
}

func (a *API) RegisterServer(ctx context.Context, grpcServer server.Server) error {
	grpcServer.RegisterServer(a.grpcServer)
	handler, prefix, err := server.CreateGateway(ctx, grpcServer, a.port)
	if err != nil {
		return err
	}
	a.RegisterHandler(prefix, handler)
	a.verifier.RegisterServer(grpcServer.AppName(), grpcServer.MethodPrefix(), grpcServer.AuthMethods())
	return nil
}

func (a *API) RegisterHandler(prefix string, handler http.Handler) {
	prefix = strings.TrimSuffix(prefix, "/")
	subRouter := a.router.PathPrefix(prefix).Name(prefix).Subrouter()
	subRouter.Use(sentryhttp.New(sentryhttp.Options{}).Handle)
	subRouter.PathPrefix("").Handler(http.StripPrefix(prefix, handler))
}

func (a *API) routeGRPC() {
	http2Route := a.router.
		MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
			return r.ProtoMajor == 2
		}).
		Subrouter()
	http2Route.
		Methods(http.MethodPost).
		Headers("Content-Type", "application/grpc").
		Handler(a.grpcServer)

	if !a.externalSecure {
		a.routeGRPCWeb(a.router)
		return
	}
	a.routeGRPCWeb(http2Route)
}

func (a *API) routeGRPCWeb(router *mux.Router) {
	router.NewRoute().
		Methods(http.MethodPost, http.MethodOptions).
		MatcherFunc(
			func(r *http.Request, _ *mux.RouteMatch) bool {
				if strings.Contains(strings.ToLower(r.Header.Get("content-type")), "application/grpc-web+") {
					return true
				}
				return strings.Contains(strings.ToLower(r.Header.Get("access-control-request-headers")), "x-grpc-web")
			}).
		Handler(
			grpcweb.WrapServer(a.grpcServer,
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
			),
		)
}

func (a *API) healthHandler() http.Handler {
	checks := []ValidationFunction{
		func(ctx context.Context) error {
			if err := a.health.Health(ctx); err != nil {
				return errors.ThrowInternal(err, "API-F24h2", "DB CONNECTION ERROR")
			}
			return nil
		},
		func(ctx context.Context) error {
			iam, err := a.health.Instance(ctx)
			if err != nil && !errors.IsNotFound(err) {
				return errors.ThrowPreconditionFailed(err, "API-dsgT2", "IAM SETUP CHECK FAILED")
			}
			if iam == nil || iam.SetupStarted < domain.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-HBfs3", "IAM NOT SET UP")
			}
			if iam.SetupDone < domain.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-DASs2", "IAM SETUP RUNNING")
			}
			return nil
		},
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", handleHealth)
	handler.HandleFunc("/ready", handleReadiness(checks))
	handler.HandleFunc("/validate", handleValidate(checks))

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
