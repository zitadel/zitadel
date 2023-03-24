package admin

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/admin/repository"
	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

const (
	adminName = "Admin-API"
)

var _ admin.AdminServiceServer = (*Server)(nil)

type Server struct {
	admin.UnimplementedAdminServiceServer
	database          string
	command           *command.Commands
	query             *query.Queries
	administrator     repository.AdministratorRepository
	assetsAPIDomain   func(context.Context) string
	userCodeAlg       crypto.EncryptionAlgorithm
	passwordHashAlg   crypto.HashAlgorithm
	auditLogRetention time.Duration
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(
	database string,
	command *command.Commands,
	query *query.Queries,
	sd systemdefaults.SystemDefaults,
	repo repository.Repository,
	externalSecure bool,
	userCodeAlg crypto.EncryptionAlgorithm,
	auditLogRetention time.Duration,
) *Server {
	return &Server{
		database:          database,
		command:           command,
		query:             query,
		administrator:     repo,
		assetsAPIDomain:   assets.AssetAPI(externalSecure),
		userCodeAlg:       userCodeAlg,
		passwordHashAlg:   crypto.NewBCrypt(sd.SecretGenerators.PasswordSaltCost),
		auditLogRetention: auditLogRetention,
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
