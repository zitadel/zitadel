package api

import (
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
)

type Config struct {
	GRPC grpc_util.Config
}

//
//func Start(ctx context.Context, conf Config, authZRepo *authz_repo.EsRepository, authZ authz.Config, defaults systemdefaults.SystemDefaults, repo repository.Repository) {
//	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig(), authZRepo, authZ, repo)
//	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())
//
//	server.StartServer(ctx, grpcServer, defaults)
//	server.StartGateway(ctx, grpcGateway)
//}
