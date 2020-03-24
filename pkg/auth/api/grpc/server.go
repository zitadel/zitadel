package grpc

import (
	grpc_utils "github.com/caos/zitadel/internal/api/grpc"
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
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_utils.ErrorHandler(),
			),
		),
	)
	RegisterAuthServiceServer(gs, s)
	return gs, nil
}
