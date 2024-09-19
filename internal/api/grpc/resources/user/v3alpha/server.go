package user

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/server"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	user "github.com/zitadel/zitadel/v2/pkg/grpc/resources/user/v3alpha"
)

var _ user.ZITADELUsersServer = (*Server)(nil)

type Server struct {
	user.UnimplementedZITADELUsersServer
	command     *command.Commands
	userCodeAlg crypto.EncryptionAlgorithm
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	userCodeAlg crypto.EncryptionAlgorithm,
) *Server {
	return &Server{
		command:     command,
		userCodeAlg: userCodeAlg,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	user.RegisterZITADELUsersServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return user.ZITADELUsers_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return user.ZITADELUsers_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return user.ZITADELUsers_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return user.RegisterZITADELUsersHandler
}
