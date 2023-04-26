package management

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

const (
	mgmtName = "Management-API"
)

var _ management.ManagementServiceServer = (*Server)(nil)

type Server struct {
	management.UnimplementedManagementServiceServer
	command           *command.Commands
	query             *query.Queries
	systemDefaults    systemdefaults.SystemDefaults
	assetAPIPrefix    func(context.Context) string
	passwordHashAlg   crypto.HashAlgorithm
	userCodeAlg       crypto.EncryptionAlgorithm
	externalSecure    bool
	auditLogRetention time.Duration
}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	sd systemdefaults.SystemDefaults,
	userCodeAlg crypto.EncryptionAlgorithm,
	externalSecure bool,
	auditLogRetention time.Duration,
) *Server {
	return &Server{
		command:           command,
		query:             query,
		systemDefaults:    sd,
		assetAPIPrefix:    assets.AssetAPI(externalSecure),
		passwordHashAlg:   crypto.NewBCrypt(sd.SecretGenerators.PasswordSaltCost),
		userCodeAlg:       userCodeAlg,
		externalSecure:    externalSecure,
		auditLogRetention: auditLogRetention,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	management.RegisterManagementServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return mgmtName
}

func (s *Server) MethodPrefix() string {
	return management.ManagementService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return management.ManagementService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return management.RegisterManagementServiceHandler
}

func (s *Server) GatewayPathPrefix() string {
	return "/management/v1"
}
