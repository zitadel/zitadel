package api

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/auth"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/api/oidc"
	auth_es "github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	authz_es "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"
	//"github.com/caos/zitadel/pkg/api/admin"
	//"github.com/caos/zitadel/pkg/api/oidc"
)

type Config struct {
	GRPC grpc.Config
	OIDC oidc.OPHandlerConfig
	Mgmt management.Config
	Auth auth.Config
	//Admin        admin.Config
}

func Start(ctx context.Context, config Config, authZ authz.Config, authZRepo *authz_es.EsRepository, sd systemdefaults.SystemDefaults, authRepo *auth_es.EsRepository, adminEnabled, managementEnabled, authEnabled, oidcEnabled bool) {
	apis := make([]server.Server, 0, 3)

	if managementEnabled {
		roles := []string{}
		managementRepo, err := mgmt_es.Start(config.Mgmt.Repository, sd, roles)
		logging.Log("API-Gd2qq").OnError(err).Fatal("error starting management repo")
		apis = append(apis, management.CreateServer(authZRepo, authZ, sd, managementRepo))
	}
	if authEnabled {
		apis = append(apis, auth.CreateServer(authZRepo, authZ, authRepo))
	}
	if adminEnabled {
		//adminRepo, err := admin_es.Start(ctx, conf.Admin, authZRepo, authZ, sd)
		//logging.Log("API-D42tq").OnError(err).Fatal("error starting auth repo")
		//admin.CreateServer(),
	}
	grpcServer := server.CreateServer(apis)
	gatewayHandler := server.CreateGatewayHandler(config.GRPC)

	for _, api := range apis {
		api.RegisterServer(grpcServer)
		gatewayHandler.RegisterGateway(ctx, api)
	}

	if oidcEnabled {
		op := oidc.NewProvider(ctx, config.OIDC, authRepo)
		gatewayHandler.RegisterHandler("/oauth/v2", op.HttpHandler().Handler)
	}

	server.Serve(ctx, grpcServer, config.GRPC.ServerPort)
	gatewayHandler.Serve(ctx)
}
