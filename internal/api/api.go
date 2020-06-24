package api

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/management"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/api/oidc"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	mgmt_es "github.com/caos/zitadel/internal/management/repository/eventsourcing"
	//"github.com/caos/zitadel/pkg/api/admin"
	//"github.com/caos/zitadel/pkg/api/auth"
	//"github.com/caos/zitadel/pkg/api/management"
	//"github.com/caos/zitadel/pkg/api/oidc"
)

type Config struct {
	GRPC grpc.Config
	OIDC oidc.Config
}

func Start(ctx context.Context, config Config, authZ authz.Config, authZRepo *authz_repo.EsRepository, sd systemdefaults.SystemDefaults, managementRepo *mgmt_es.EsRepository) {
	apis := make([]server.Server, 0, 3)
	apis = append(apis, management.CreateServer(authZRepo, authZ, sd, managementRepo))
	//	admin.CreateServer(),
	//	auth.CreateServer(),
	//}

	grpcServer := server.CreateServer(apis)
	gatewayHandler := server.CreateGatewayHandler(config.GRPC)

	for _, api := range apis {
		api.RegisterServer(grpcServer)
		gatewayHandler.RegisterGateway(ctx, api)
	}
	gatewayHandler.RegisterHandler("/oauth/v2", nil)

	server.Serve(ctx, grpcServer, config.GRPC.ServerPort)
	gatewayHandler.Serve(ctx)
}
