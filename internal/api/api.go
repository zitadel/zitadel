package api

import (
	"context"
	"net/http"

	"github.com/caos/logging"
	admin_es "github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	authz_es "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/caos/zitadel/internal/telemetry/metrics/otel"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	view_model "github.com/caos/zitadel/internal/view/model"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

type Config struct {
	GRPC grpc_util.Config
	OIDC oidc.OPHandlerConfig
}

type API struct {
	grpcServer     *grpc.Server
	gatewayHandler *server.GatewayHandler
	verifier       *authz.TokenVerifier
	serverPort     string
	health         health
	auth           auth
	admin          admin
}

type health interface {
	Health(ctx context.Context) error
	IamByID(ctx context.Context) (*iam_model.IAM, error)
	VerifierClientID(ctx context.Context, appName string) (string, error)
}

type auth interface {
	ActiveUserSessionCount() int64
}

type admin interface {
	GetViews() ([]*view_model.View, error)
	GetSpoolerDiv(database, viewName string) int64
}

func Create(config Config, authZ authz.Config, authZRepo *authz_es.EsRepository, authRepo *auth_es.EsRepository, adminRepo *admin_es.EsRepository, sd systemdefaults.SystemDefaults) *API {
	api := &API{
		serverPort: config.GRPC.ServerPort,
	}
	api.verifier = authz.Start(authZRepo)
	api.health = authZRepo
	api.auth = authRepo
	api.admin = adminRepo
	api.grpcServer = server.CreateServer(api.verifier, authZ, sd.DefaultLanguage)
	api.gatewayHandler = server.CreateGatewayHandler(config.GRPC)
	api.RegisterHandler("", api.healthHandler())

	return api
}

func (a *API) RegisterServer(ctx context.Context, server server.Server) {
	server.RegisterServer(a.grpcServer)
	a.gatewayHandler.RegisterGateway(ctx, server)
	a.verifier.RegisterServer(server.AppName(), server.MethodPrefix(), server.AuthMethods())
}

func (a *API) RegisterHandler(prefix string, handler http.Handler) {
	a.gatewayHandler.RegisterHandler(prefix, handler)
}

func (a *API) Start(ctx context.Context) {
	server.Serve(ctx, a.grpcServer, a.serverPort)
	a.gatewayHandler.Serve(ctx)
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
			iam, err := a.health.IamByID(ctx)
			if err != nil && !errors.IsNotFound(err) {
				return errors.ThrowPreconditionFailed(err, "API-dsgT2", "IAM SETUP CHECK FAILED")
			}
			if iam == nil || iam.SetUpStarted < iam_model.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-HBfs3", "IAM NOT SET UP")
			}
			if iam.SetUpDone < iam_model.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-DASs2", "IAM SETUP RUNNING")
			}
			return nil
		},
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", handleHealth)
	handler.HandleFunc("/ready", handleReadiness(checks))
	handler.HandleFunc("/validate", handleValidate(checks))
	handler.HandleFunc("/clientID", a.handleClientID)
	handler.Handle("/metrics", a.handleMetrics())

	return handler
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	logging.Log("API-Hfss2").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(r.Context())).Error("error writing ok for health")
}

func handleReadiness(checks []ValidationFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errors := validate(r.Context(), checks)
		if len(errors) == 0 {
			http_util.MarshalJSON(w, "ok", nil, http.StatusOK)
			return
		}
		http_util.MarshalJSON(w, nil, errors[0], http.StatusPreconditionFailed)
	}
}

func handleValidate(checks []ValidationFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errors := validate(r.Context(), checks)
		if len(errors) == 0 {
			http_util.MarshalJSON(w, "ok", nil, http.StatusOK)
			return
		}
		http_util.MarshalJSON(w, errors, nil, http.StatusOK)
	}
}

func (a *API) handleClientID(w http.ResponseWriter, r *http.Request) {
	id, err := a.health.VerifierClientID(r.Context(), "Zitadel Console")
	if err != nil {
		http_util.MarshalJSON(w, nil, err, http.StatusPreconditionFailed)
		return
	}
	http_util.MarshalJSON(w, id, nil, http.StatusOK)
}

func (a *API) handleMetrics() http.Handler {
	a.registerActiveSessionCounters()
	a.registerSpoolerDivCounters()
	return metrics.GetExporter()
}

func (a *API) registerActiveSessionCounters() {
	metrics.RegisterValueObserver(
		metrics.ActiveSessionCounter,
		metrics.ActiveSessionCounterDescription,
		func(ctx context.Context, result metric.Int64ObserverResult) {
			result.Observe(
				a.auth.ActiveUserSessionCount(),
			)
		},
	)
}

func (a *API) registerSpoolerDivCounters() {
	views, err := a.admin.GetViews()
	if err != nil {
		logging.Log("API-3M8sd").WithError(err).Error("could not read views for metrics")
		return
	}
	metrics.RegisterValueObserver(
		metrics.SpoolerDivCounter,
		metrics.SpoolerDivCounterDescription,
		func(ctx context.Context, result metric.Int64ObserverResult) {
			for _, view := range views {
				labels := map[string]interface{}{
					metrics.Database: view.Database,
					metrics.ViewName: view.ViewName,
				}
				result.Observe(
					a.admin.GetSpoolerDiv(view.Database, view.ViewName),
					otel.MapToKeyValue(labels)...,
				)
			}
		},
	)
}

type ValidationFunction func(ctx context.Context) error

func validate(ctx context.Context, validations []ValidationFunction) []error {
	errors := make([]error, 0)
	for _, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Error("validation failed")
			errors = append(errors, err)
		}
	}
	return errors
}
