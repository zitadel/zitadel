package api

import (
	"context"
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	server2 "github.com/caos/zitadel/v2/internal/api/grpc/server"
	"github.com/caos/zitadel/v2/internal/api/oidc"

	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
)

type API struct {
	port           string
	grpcServer     *grpc.Server
	gatewayHandler *server2.GatewayHandler
	verifier       *internal_authz.TokenVerifier
	router         *mux.Router
}

type Config struct {
	oidc.OPHandlerConfig
}

func New(ctx context.Context, port string, router *mux.Router, verifier *internal_authz.TokenVerifier, authZ internal_authz.Config, sd systemdefaults.SystemDefaults) *API {
	api := &API{
		port:     port,
		verifier: verifier,
		router:   router,
	}
	api.grpcServer = server2.CreateServer(api.verifier, authZ, sd.DefaultLanguage)
	api.gatewayHandler = server2.CreateGatewayHandler(port)

	return api
}

func (a *API) RegisterServer(ctx context.Context, server server2.Server) {
	server.RegisterServer(a.grpcServer)
	a.gatewayHandler.RegisterGateway(ctx, server)
	a.verifier.RegisterServer(server.AppName(), server.MethodPrefix(), server.AuthMethods())
}

func (a *API) RegisterHandler(prefix string, handler http.Handler) {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	subRouter := a.router.PathPrefix(prefix).Subrouter()
	subRouter.PathPrefix("/").Handler(http.StripPrefix(prefix, sentryHandler.Handle(handler)))
}
func (a *API) Router() *mux.Router {
	a.routeGRPC()
	a.routeHTTP()
	return a.router
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

func (a *API) routeHTTP() {
	a.router.PathPrefix("/").Handler(a.gatewayHandler.Router())
	//http1Router.PathPrefix("/").Handler(http.StripPrefix("/", a.gatewayHandler.Router()))
}
