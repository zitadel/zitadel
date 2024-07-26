package settings

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

var _ settings.SettingsServiceServer = (*Server)(nil)

type Server struct {
	settings.UnimplementedSettingsServiceServer
	command         *command.Commands
	query           *query.Queries
	assetsAPIDomain func(context.Context) string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	externalSecure bool,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		assetsAPIDomain: assets.AssetAPI(externalSecure),
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	settings.RegisterSettingsServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return settings.SettingsService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return settings.SettingsService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return settings.SettingsService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return settings.RegisterSettingsServiceHandler
}
