package auth

import (
	"context"
	"net/http"

	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	service auth.AuthServiceServer
}

func New() *Server {
	return &Server{
		service: auth.UnimplementedAuthServiceServer{},
	}
}

func (s *Server) RegisterGRPC(srv *grpc.Server) {
	auth.RegisterAuthServiceServer(srv, s.service)
}

func (s *Server) RegisterRESTGateway(ctx context.Context, m *http.ServeMux, grpcMux *runtime.ServeMux) error {
	conn, err := grpc.Dial(":50002", grpc.WithInsecure())
	if err != nil {
		return err
	}

	m.Handle("/api/auth/v1", grpcMux)

	return auth.RegisterAuthServiceHandler(ctx, grpcMux, conn)
}

func (s *Server) registerGRPCWebGateway() {}
