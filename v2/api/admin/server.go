package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	service admin.AdminServiceServer
}

func New(adminSvc admin.AdminServiceServer) *Server {
	return &Server{
		service: adminSvc,
	}
}

func (s *Server) RegisterGRPC(srv *grpc.Server) {
	admin.RegisterAdminServiceServer(srv, s.service)
}

func (s *Server) RegisterRESTGateway(ctx context.Context, grpcMux *runtime.ServeMux) error {
	conn, err := grpc.Dial(":50002", grpc.WithInsecure())
	if err != nil {
		return err
	}

	return admin.RegisterAdminServiceHandler(ctx, grpcMux, conn)
}

func (s *Server) ServicePrefix() string {
	return "/admin/v1"
}

func (s *Server) AppName() string {
	return "Admin-API"
}

func (s *Server) MethodPrefix() string {
	return admin.AdminService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return admin.AdminService_AuthMethods
}
