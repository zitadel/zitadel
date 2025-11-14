package auth

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
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
	assetsAPIDomain func(context.Context) string
	userCodeAlg     crypto.EncryptionAlgorithm
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands,
	query *query.Queries,
	authRepo repository.Repository,
	defaults systemdefaults.SystemDefaults,
	userCodeAlg crypto.EncryptionAlgorithm,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		repo:            authRepo,
		defaults:        defaults,
		assetsAPIDomain: assets.AssetAPI(),
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
	return auth.AuthService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return auth.AuthService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return auth.RegisterAuthServiceHandler
}

func (s *Server) GatewayPathPrefix() string {
	return GatewayPathPrefix()
}

func GatewayPathPrefix() string {
	return "/auth/v1"
}
