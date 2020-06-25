package auth

import (
	"google.golang.org/grpc"

	auth_util "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/auth/auth"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	auth_grpc "github.com/caos/zitadel/pkg/auth/grpc"
)

var _ auth_grpc.AuthServiceServer = (*Server)(nil)

type Server struct {
	repo     repository.Repository
	verifier *auth.TokenVerifier
	authZ    auth_util.Config
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(authZRepo *authz_repo.EsRepository, authZ auth_util.Config, authRepo repository.Repository) *Server {
	return &Server{
		repo:     authRepo,
		authZ:    authZ,
		verifier: auth.Start(authZRepo),
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	auth_grpc.RegisterAuthServiceServer(grpcServer, s)
}

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return auth_grpc.AuthService_Authorization_Interceptor(nil, nil)
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return auth_grpc.RegisterAuthServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/auth/v1"
}
