package grpc

import (
	"github.com/caos/zitadel/internal/api/auth"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	mgmt_auth "github.com/caos/zitadel/internal/management/auth"
	"github.com/caos/zitadel/internal/management/repository"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var _ ManagementServiceServer = (*Server)(nil)

type Server struct {
	port           string
	project        repository.ProjectRepository
	policy         repository.PolicyRepository
	org            repository.OrgRepository
	user           repository.UserRepository
	usergrant      repository.UserGrantRepository
	iam            repository.IamRepository
	verifier       *mgmt_auth.TokenVerifier
	authZ          auth.Config
	systemDefaults systemdefaults.SystemDefaults
}

func StartServer(conf grpc_util.ServerConfig, authZRepo *authz_repo.EsRepository, authZ auth.Config, sd systemdefaults.SystemDefaults, repo repository.Repository) *Server {
	return &Server{
		port:           conf.Port,
		project:        repo,
		policy:         repo,
		org:            repo,
		user:           repo,
		usergrant:      repo,
		iam:            repo,
		authZ:          authZ,
		verifier:       mgmt_auth.Start(authZRepo),
		systemDefaults: sd,
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
				ManagementService_Authorization_Interceptor(s.verifier, &s.authZ),
			),
		),
	)
	RegisterManagementServiceServer(gs, s)
	return gs, nil
}
