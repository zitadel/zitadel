package admin

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	admin_grpc "github.com/caos/zitadel/pkg/admin/grpc"
)

const (
	adminName = "Admin-API"
)

var _ admin_grpc.AdminServiceServer = (*Server)(nil)

type Server struct {
	org           repository.OrgRepository
	iam           repository.IamRepository
	administrator repository.AdministratorRepository
	repo          repository.Repository
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(repo repository.Repository) *Server {
	return &Server{
		org:           repo,
		iam:           repo,
		administrator: repo,
		repo:          repo,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	admin_grpc.RegisterAdminServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return adminName
}

func (s *Server) MethodPrefix() string {
	return admin_grpc.AdminService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return admin_grpc.AdminService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return admin_grpc.RegisterAdminServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/admin/v1"
}
