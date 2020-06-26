package admin

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	admin_grpc "github.com/caos/zitadel/pkg/admin/grpc"
)

const (
	adminName = "Admin-API"
)

var _ admin_grpc.AdminServiceServer = (*Server)(nil)

type Server struct {
	port string
	org  repository.OrgRepository
	//verifier      authz.TokenVerifier
	authZ         authz.Config
	iam           repository.IamRepository
	administrator repository.AdministratorRepository
	repo          repository.Repository
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(authZRepo *authz_repo.EsRepository, authZ authz.Config, repo repository.Repository) *Server {
	return &Server{
		org:           repo,
		iam:           repo,
		administrator: repo,
		repo:          repo,
		authZ:         authZ,
		//verifier:      admin_auth.Start(authZRepo),
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	admin_grpc.RegisterAdminServiceServer(grpcServer, s)
}

//func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
//	return admin_grpc.AdminService_Authorization_Interceptor(nil, nil)
//}

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
