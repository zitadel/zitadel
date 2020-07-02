package api

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/api/oidc"
	authz_es "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
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
}

func Create(config Config, authZ authz.Config, authZRepo *authz_es.EsRepository) *API {
	apis := make([]server.Server, 0, 3)

	api := &API{
		serverPort: config.GRPC.ServerPort,
	}
	api.verifier = authz.Start(authZRepo)
	api.grpcServer = server.CreateServer(apis, api.verifier, authZ)
	api.gatewayHandler = server.CreateGatewayHandler(config.GRPC)

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
