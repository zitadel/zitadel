package grpc

import (
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ ManagementServiceServer = (*Server)(nil)

type Server struct {
	port string
}

func StartServer(conf *grpc_util.ServerConfig) *Server {
	return &Server{
		port: conf.Port,
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
