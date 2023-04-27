package user

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
)

var _ UserServiceServer = (*Server)(nil)

type Server struct {
	UnimplementedUserServiceServer
	command     *command.Commands
	query       *query.Queries
	userCodeAlg crypto.EncryptionAlgorithm
}

type Config struct{}

func CreateServer(command *command.Commands, query *query.Queries, userCodeAlg crypto.EncryptionAlgorithm) *Server {
	return &Server{
		command:     command,
		query:       query,
		userCodeAlg: userCodeAlg,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	RegisterUserServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return UserService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return UserService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return UserService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return RegisterUserServiceHandler
}
