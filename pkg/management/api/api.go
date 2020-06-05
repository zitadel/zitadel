package api

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/management/repository"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/management/api/grpc"
)

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf Config, authZRepo *authz_repo.EsRepository, authZ auth.Config, repo repository.Repository) {
	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig(), authZRepo, authZ, repo)
	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())

	server.StartServer(ctx, grpcServer)
	server.StartGateway(ctx, grpcGateway)
}
