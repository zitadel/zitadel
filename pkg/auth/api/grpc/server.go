package grpc

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ AuthServiceServer = (*Server)(nil)

type Config struct {
	Port        string
	SearchLimit int
}

type Server struct {
	port        string
	searchLimit int
}

func StartServer(conf Config) *Server {
	return &Server{
		port:        conf.Port,
		searchLimit: conf.SearchLimit,
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
