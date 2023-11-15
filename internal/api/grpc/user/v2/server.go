package user

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/server"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/query"
	user "github.com/zitadel/zitadel/v2/pkg/grpc/user/v2beta"
)

var _ user.UserServiceServer = (*Server)(nil)

type Server struct {
	user.UnimplementedUserServiceServer
	command     *command.Commands
	query       *query.Queries
	userCodeAlg crypto.EncryptionAlgorithm
	idpAlg      crypto.EncryptionAlgorithm
	idpCallback func(ctx context.Context) string
	samlRootURL func(ctx context.Context, idpID string) string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	idpAlg crypto.EncryptionAlgorithm,
	idpCallback func(ctx context.Context) string,
	samlRootURL func(ctx context.Context, idpID string) string,
) *Server {
	return &Server{
		command:     command,
		query:       query,
		userCodeAlg: userCodeAlg,
		idpAlg:      idpAlg,
		idpCallback: idpCallback,
		samlRootURL: samlRootURL,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	user.RegisterUserServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return user.UserService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return user.UserService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return user.UserService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return user.RegisterUserServiceHandler
}
