package auth

import (
	"github.com/caos/logging"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/crypto"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

var _ auth.AuthServiceServer = (*Server)(nil)

const (
	authName = "Auth-API"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	command         *command.Commands
	query           *query.Queries
	repo            repository.Repository
	defaults        systemdefaults.SystemDefaults
	assetsAPIDomain string
	userCodeAlg     crypto.EncryptionAlgorithm
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands,
	query *query.Queries,
	authRepo repository.Repository,
	defaults systemdefaults.SystemDefaults,
	assetsAPIDomain string,
	keyStorage crypto.KeyStorage,
	userEncryptionConfig *crypto.KeyConfig,
) *Server {
	userCodeAlg, err := crypto.NewAESCrypto(userEncryptionConfig, keyStorage)
	logging.OnError(err).Fatal("unable to initialise user code algorithm")
	return &Server{
		command:         command,
		query:           query,
		repo:            authRepo,
		defaults:        defaults,
		assetsAPIDomain: assetsAPIDomain,
		userCodeAlg:     userCodeAlg,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	auth.RegisterAuthServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return authName
}

func (s *Server) MethodPrefix() string {
	return auth.AuthService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return auth.AuthService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return auth.RegisterAuthServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/auth/v1"
}
