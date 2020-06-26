package auth

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	auth_grpc "github.com/caos/zitadel/pkg/auth/grpc"
)

var _ auth_grpc.AuthServiceServer = (*Server)(nil)

const (
	authName = "Auth-API"
)

type Server struct {
	repo repository.Repository
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(authRepo repository.Repository) *Server {
	return &Server{
		repo: authRepo,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	auth_grpc.RegisterAuthServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return authName
}

func (s *Server) MethodPrefix() string {
	return auth_grpc.AuthService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return auth_grpc.AuthService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return auth_grpc.RegisterAuthServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/auth/v1"
}
