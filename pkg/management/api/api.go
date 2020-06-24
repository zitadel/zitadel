package api

import (
	"context"

	"github.com/caos/zitadel/pkg/management/api/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/management/repository"
)

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf Config, authZRepo *authz_repo.EsRepository, authZ authz.Config, defaults systemdefaults.SystemDefaults, repo repository.Repository) {
	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig(), authZRepo, authZ, defaults, repo)
	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())

	server.StartServer(ctx, grpcServer, defaults)
	server.StartGateway(ctx, grpcGateway)
}
