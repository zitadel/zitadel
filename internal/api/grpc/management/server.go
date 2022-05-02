package management

import (
	"time"

	"google.golang.org/grpc"

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
	assetAPIPrefix    string
	passwordHashAlg   crypto.HashAlgorithm
	userCodeAlg       crypto.EncryptionAlgorithm
	externalSecure    bool
	issuerPath        string
	auditLogRetention time.Duration
}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	sd systemdefaults.SystemDefaults,
	assetAPIPrefix string,
	userCodeAlg crypto.EncryptionAlgorithm,
	externalSecure bool,
	issuerPath string,
	auditLogRetention time.Duration,
) *Server {
	return &Server{
		command:           command,
		query:             query,
		systemDefaults:    sd,
		assetAPIPrefix:    assetAPIPrefix,
		passwordHashAlg:   crypto.NewBCrypt(sd.SecretGenerators.PasswordSaltCost),
		userCodeAlg:       userCodeAlg,
		externalSecure:    externalSecure,
		issuerPath:        issuerPath,
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
	return management.ManagementService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return management.ManagementService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return management.RegisterManagementServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/management/v1"
}
