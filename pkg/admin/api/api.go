package api

import (
	"context"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/admin/api/grpc"
)

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf Config) {
	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig())
	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())

	server.StartServer(ctx, grpcServer)
	server.StartGateway(ctx, grpcGateway)
}
