package webkey

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

type Server struct {
	v3alpha.UnimplementedWebKeyServiceServer
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
	v3alpha.RegisterWebKeyServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return v3alpha.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return v3alpha.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return v3alpha.WebKeyService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return v3alpha.RegisterWebKeyServiceHandler
}
