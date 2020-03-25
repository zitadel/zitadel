package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

var _ AuthServiceServer = (*Server)(nil)

type Server struct {
	port        string
	searchLimit int
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
	RegisterAuthServiceServer(gs, s)
	return gs, nil
}
