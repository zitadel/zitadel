package webkey

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

type Server struct {
	webkey.UnimplementedWebKeyServiceServer
	command *command.Commands
	query   *query.Queries
}

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
	webkey.RegisterWebKeyServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return webkey.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return webkey.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return webkey.WebKeyService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return webkey.RegisterWebKeyServiceHandler
}
