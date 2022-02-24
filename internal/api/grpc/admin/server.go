package admin

import (
	"github.com/caos/logging"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/crypto"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

const (
	adminName = "Admin-API"
)

var _ admin.AdminServiceServer = (*Server)(nil)

type Server struct {
	admin.UnimplementedAdminServiceServer
	command         *command.Commands
	query           *query.Queries
	administrator   repository.AdministratorRepository
	iamDomain       string
	assetsAPIDomain string

	UserCodeAlg crypto.EncryptionAlgorithm
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	iamDomain,
	assetsAPIDomain string,
	keyStorage crypto.KeyStorage,
	userEncryptionConfig *crypto.KeyConfig,
) *Server {
	userCodeAlg, err := crypto.NewAESCrypto(userEncryptionConfig, keyStorage)
	logging.OnError(err).Fatal("unable to initialise user code algorithm")
	return &Server{
		command:         command,
		query:           query,
		administrator:   repo,
		iamDomain:       iamDomain,
		assetsAPIDomain: assetsAPIDomain,
		UserCodeAlg:     userCodeAlg,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	admin.RegisterAdminServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return adminName
}

func (s *Server) MethodPrefix() string {
	return admin.AdminService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return admin.AdminService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return admin.RegisterAdminServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/admin/v1"
}
