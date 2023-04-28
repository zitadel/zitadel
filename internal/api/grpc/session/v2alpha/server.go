package session

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
)

var _ SessionServiceServer = (*Server)(nil)

type Server struct {
	UnimplementedSessionServiceServer
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
	RegisterSessionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return SessionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return SessionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return SessionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return RegisterSessionServiceHandler
}
