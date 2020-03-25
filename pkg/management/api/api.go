package api

import (
	"context"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/management/api/grpc"
)

type API struct {
	grpcServer grpc.Server
	gateway    grpc.Gateway
}

type Config struct {
	GRPCServer grpc.Config
	Gateway    grpc.GatewayConfig
}

func Start(ctx context.Context, conf *Config) error {
	api := &API{
		grpcServer: *grpc.StartServer(conf.GRPCServer),
		gateway:    *grpc.StartGateway(conf.Gateway),
	}
	server.StartServer(ctx, &api.grpcServer)
	server.StartGateway(ctx, &api.gateway)

	return nil
}
