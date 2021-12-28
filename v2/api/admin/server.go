package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	service admin.AdminServiceServer
}

func New() *Server {
	return &Server{
		service: admin.UnimplementedAdminServiceServer{},
	}
}

func (s *Server) RegisterGRPC(srv *grpc.Server) {
	admin.RegisterAdminServiceServer(srv, s.service)
}

func (s *Server) RegisterRESTGateway(ctx context.Context, m *runtime.ServeMux) error {
	conn, err := grpc.Dial(":50002", grpc.WithInsecure())
	if err != nil {
		return err
	}
	return admin.RegisterAdminServiceHandler(ctx, m, conn)
}

func (s *Server) registerGRPCWebGateway() {}
