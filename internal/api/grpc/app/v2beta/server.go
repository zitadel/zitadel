package app

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

var _ app.AppServiceServer = (*Server)(nil)

type Server struct {
	app.UnimplementedAppServiceServer
	command         *command.Commands
	query           *query.Queries
	systemDefaults  systemdefaults.SystemDefaults
	checkPermission domain.PermissionCheck
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		checkPermission: checkPermission,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	app.RegisterAppServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return app.AppService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return app.AppService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return app.AppService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return app.RegisterAppServiceHandler
}
