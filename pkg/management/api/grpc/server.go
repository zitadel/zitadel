package grpc

import (
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	"github.com/caos/zitadel/internal/management/repository"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ ManagementServiceServer = (*Server)(nil)

type Server struct {
	port    string
	project repository.ProjectRepository
}

func StartServer(conf grpc_util.ServerConfig, repo repository.Repository) *Server {
	return &Server{
		port:    conf.Port,
		project: repo,
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
	RegisterManagementServiceServer(gs, s)
	return gs, nil
}
