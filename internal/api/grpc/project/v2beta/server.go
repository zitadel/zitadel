package project

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

var _ project.ProjectServiceServer = (*Server)(nil)

type Server struct {
	project.UnimplementedProjectServiceServer
	systemDefaults systemdefaults.SystemDefaults
	command        *command.Commands
	query          *query.Queries
}

type Config struct{}

func CreateServer(
	systemDefaults systemdefaults.SystemDefaults,
	command *command.Commands,
	query *query.Queries,
) *Server {
	return &Server{
		systemDefaults: systemDefaults,
		command:        command,
		query:          query,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	project.RegisterProjectServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return project.ProjectService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return project.ProjectService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return project.ProjectService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return project.RegisterProjectServiceHandler
}
