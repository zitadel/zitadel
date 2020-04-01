package api

import (
	"context"
	"github.com/caos/zitadel/internal/management/repository"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/management/api/grpc"
)

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf Config, repo repository.Repository) {
	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig(), repo)
	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())

	server.StartServer(ctx, grpcServer)
	server.StartGateway(ctx, grpcGateway)
}
