package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/caos/logging"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/authz/repository"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type API struct {
	port       string
	grpcServer *grpc.Server
	verifier   *internal_authz.TokenVerifier
	health     health
	router     *mux.Router
}

type health interface {
	Health(ctx context.Context) error
	IAMByID(ctx context.Context, id string) (*query.IAM, error)
	VerifierClientID(ctx context.Context, appName string) (string, string, error)
}

func New(
	port string,
	router *mux.Router,
	repo *struct {
		repository.Repository
		*query.Queries
	},
	authZ internal_authz.Config,
	sd systemdefaults.SystemDefaults,
) *API {
	verifier := internal_authz.Start(repo)
	api := &API{
		port:     port,
		verifier: verifier,
		health:   repo,
		router:   router,
	}
	api.grpcServer = server.CreateServer(api.verifier, authZ, sd.DefaultLanguage)
	api.routeGRPC()

	api.RegisterHandler("", api.healthHandler()) //TODO: do we need a prefix?

	return api
}

func (a *API) RegisterServer(ctx context.Context, grpcServer server.Server) {
	grpcServer.RegisterServer(a.grpcServer)
	handler, prefix := server.CreateGateway(ctx, grpcServer, a.port)
	a.RegisterHandler(prefix, handler)
	a.verifier.RegisterServer(grpcServer.AppName(), grpcServer.MethodPrefix(), grpcServer.AuthMethods())
}

func (a *API) RegisterHandler(prefix string, handler http.Handler) {
	prefix = strings.TrimSuffix(prefix, "/")
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	subRouter := a.router.PathPrefix(prefix).Subrouter()
	subRouter.PathPrefix("/").Handler(http.StripPrefix(prefix, sentryHandler.Handle(handler)))
}

func (a *API) routeGRPC() {
	http2Route := a.router.Methods(http.MethodPost). //TODO: grpc-web is called with http/1.1
								MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
			return r.ProtoMajor == 2
		}).
		Subrouter()
	http2Route.Headers("Content-Type", "application/grpc").Handler(a.grpcServer)
	a.router.NewRoute().HeadersRegexp("Content-Type", "application/grpc-web.*").Handler(grpcweb.WrapServer(a.grpcServer))
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
			iam, err := a.health.IAMByID(ctx, domain.IAMID)
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
	handler.HandleFunc("/clientID", a.handleClientID) //TODO: remove?

	return handler
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	logging.Log("API-Hfss2").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(r.Context())).Error("error writing ok for health")
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

func (a *API) handleClientID(w http.ResponseWriter, r *http.Request) {
	id, _, err := a.health.VerifierClientID(r.Context(), "Zitadel Console")
	if err != nil {
		http_util.MarshalJSON(w, nil, err, http.StatusPreconditionFailed)
		return
	}
	http_util.MarshalJSON(w, id, nil, http.StatusOK)
}

type ValidationFunction func(ctx context.Context) error

func validate(ctx context.Context, validations []ValidationFunction) []error {
	errs := make([]error, 0)
	for _, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Error("validation failed")
			errs = append(errs, err)
		}
	}
	return errs
}
