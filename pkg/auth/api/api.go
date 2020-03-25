package api

import (
	"context"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/auth/api/grpc"
)

type API struct {
	grpcServer grpc.Server
	gateway    grpc.Gateway
}

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf *Config) error {
	api := &API{
		grpcServer: *grpc.StartServer(conf.GRPC.ToServerConfig()),
		gateway:    *grpc.StartGateway(conf.GRPC.ToGatewayConfig()),
	}
	server.StartServer(ctx, &api.grpcServer)
	server.StartGateway(ctx, &api.gateway)

	return nil
}
