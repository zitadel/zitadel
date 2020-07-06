package admin

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

const (
	adminName = "Admin-API"
)

var _ admin.AdminServiceServer = (*Server)(nil)

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
	admin.RegisterAdminServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return adminName
}

func (s *Server) MethodPrefix() string {
	return admin.AdminService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return admin.AdminService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return admin.RegisterAdminServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/admin/v1"
}
