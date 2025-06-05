package internal_permission

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

var _ internal_permission.InternalPermissionServiceServer = (*Server)(nil)

type Server struct {
	internal_permission.UnimplementedInternalPermissionServiceServer
	command         *command.Commands
	query           *query.Queries
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
	internal_permission.RegisterInternalPermissionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return internal_permission.InternalPermissionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return internal_permission.InternalPermissionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return internal_permission.InternalPermissionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return internal_permission.RegisterInternalPermissionServiceHandler
}
