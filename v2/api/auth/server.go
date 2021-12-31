package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	service auth.AuthServiceServer
}

func New(authSvc auth.AuthServiceServer) *Server {
	return &Server{
		service: authSvc,
	}
}

func (s *Server) RegisterGRPC(srv *grpc.Server) {
	auth.RegisterAuthServiceServer(srv, s.service)
}

func (s *Server) RegisterRESTGateway(ctx context.Context, grpcMux *runtime.ServeMux) error {
	conn, err := grpc.Dial(":50002", grpc.WithInsecure())
	if err != nil {
		return err
	}

	return auth.RegisterAuthServiceHandler(ctx, grpcMux, conn)
}

func (s *Server) ServicePrefix() string {
	return "/auth/v1"
}

func (s *Server) AppName() string {
	return "Auth-API"
}

func (s *Server) MethodPrefix() string {
	return auth.AuthService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return auth.AuthService_AuthMethods
}
