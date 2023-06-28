package oidc

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
)

var _ oidc_pb.OIDCServiceServer = (*Server)(nil)

type Server struct {
	oidc_pb.UnimplementedOIDCServiceServer
	command *command.Commands
	query   *query.Queries
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
) *Server {
	return &Server{
		command: command,
		query:   query,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	oidc_pb.RegisterOIDCServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return oidc_pb.OIDCService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return oidc_pb.OIDCService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return map[string]authz.Option{}
	//return oidc_pb.OIDCService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return nil
	//return oidc_pb.RegisterOIDCServiceHandler
}
