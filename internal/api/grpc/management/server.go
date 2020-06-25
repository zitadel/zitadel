package management

import (
	"context"
	"errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	mgmt_auth "github.com/caos/zitadel/internal/management/auth"
	"github.com/caos/zitadel/internal/management/repository"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	mgmt_grpc "github.com/caos/zitadel/pkg/management/grpc"
)

var _ mgmt_grpc.ManagementServiceServer = (*Server)(nil)

type Server struct {
	project        repository.ProjectRepository
	policy         repository.PolicyRepository
	org            repository.OrgRepository
	user           repository.UserRepository
	usergrant      repository.UserGrantRepository
	iam            repository.IamRepository
	verifier       *mgmt_auth.TokenVerifier
	authZ          authz.Config
	systemDefaults systemdefaults.SystemDefaults
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(authZRepo *authz_repo.EsRepository, authZ authz.Config, sd systemdefaults.SystemDefaults, repo repository.Repository) *Server {
	return &Server{
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

func (s *Server) GRPCServer(defaults systemdefaults.SystemDefaults) (*grpc.Server, error) {
	gs := grpc.NewServer(
		middleware.TracingStatsServer("/Healthz", "/Ready", "/Validate"),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.ErrorHandler(defaults.DefaultLanguage),
				mgmt_grpc.ManagementService_Authorization_Interceptor(s.verifier, &s.authZ),
			),
		),
	)
	mgmt_grpc.RegisterManagementServiceServer(gs, s)
	return gs, nil
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	mgmt_grpc.RegisterManagementServiceServer(grpcServer, s)
}

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return mgmt_grpc.ManagementService_Authorization_Interceptor(nil, nil)
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return mgmt_grpc.RegisterManagementServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/mgmt/v1"
}

func (s *Server) Validations() map[string]server.ValidationFunction {
	return map[string]server.ValidationFunction{
		"Test": func(_ context.Context) error { return errors.New("Test") },
	}
}
