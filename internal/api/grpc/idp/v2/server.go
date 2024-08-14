package idp

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/idp/v2"
)

var _ idp.IdentityProviderServiceServer = (*Server)(nil)

type Server struct {
	idp.UnimplementedIdentityProviderServiceServer
	command *command.Commands
	query   *query.Queries

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
	idp.RegisterIdentityProviderServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return idp.IdentityProviderService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return idp.IdentityProviderService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return idp.IdentityProviderService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return idp.RegisterIdentityProviderServiceHandler
}
