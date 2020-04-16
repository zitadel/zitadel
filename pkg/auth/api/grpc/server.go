package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	auth_util "github.com/caos/zitadel/internal/api/auth"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	"github.com/caos/zitadel/internal/auth/auth"
	"github.com/caos/zitadel/internal/auth/repository"
)

var _ AuthServiceServer = (*Server)(nil)

type Server struct {
	port     string
	repo     repository.Repository
	verifier *auth.TokenVerifier
	authZ    auth_util.Config
}

func StartServer(conf grpc_util.ServerConfig, authZ auth_util.Config, repo repository.Repository) *Server {
	return &Server{
		port:     conf.Port,
		repo:     repo,
		authZ:    authZ,
		verifier: auth.Start(),
	}
}

func (s *Server) GRPCPort() string {
	return s.port
}

func (s *Server) GRPCServer() (*grpc.Server, error) {
	gs := grpc.NewServer(
		middleware.TracingStatsServer("/Healthz", "/Ready", "/Validate"),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.ErrorHandler(),
			),
		),
	)
	RegisterAuthServiceServer(gs, s)
	return gs, nil
}
