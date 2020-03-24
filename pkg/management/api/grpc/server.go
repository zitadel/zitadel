package grpc

import (
	"google.golang.org/grpc"
)

var _ ManagementServiceServer = (*Server)(nil)

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
	gs := grpc.NewServer()
	RegisterManagementServiceServer(gs, s)
	return gs, nil
}
