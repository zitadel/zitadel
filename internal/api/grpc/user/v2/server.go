package user

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
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

	assetAPIPrefix func(context.Context) string

	checkPermission domain.PermissionCheck
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	idpAlg crypto.EncryptionAlgorithm,
	idpCallback func(ctx context.Context) string,
	samlRootURL func(ctx context.Context, idpID string) string,
	assetAPIPrefix func(ctx context.Context) string,
	checkPermission domain.PermissionCheck,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		userCodeAlg:     userCodeAlg,
		idpAlg:          idpAlg,
		idpCallback:     idpCallback,
		samlRootURL:     samlRootURL,
		assetAPIPrefix:  assetAPIPrefix,
		checkPermission: checkPermission,
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
