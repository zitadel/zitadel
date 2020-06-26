package management

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/management/repository"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	mgmt_grpc "github.com/caos/zitadel/pkg/management/grpc"
)

const (
	mgmtName = "Management-API"
)

var _ mgmt_grpc.ManagementServiceServer = (*Server)(nil)

type Server struct {
	project        repository.ProjectRepository
	policy         repository.PolicyRepository
	org            repository.OrgRepository
	user           repository.UserRepository
	usergrant      repository.UserGrantRepository
	iam            repository.IamRepository
	authZ          authz.Config
	systemDefaults systemdefaults.SystemDefaults
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(repo repository.Repository, sd systemdefaults.SystemDefaults) *Server {
	return &Server{
		project:        repo,
		policy:         repo,
		org:            repo,
		user:           repo,
		usergrant:      repo,
		iam:            repo,
		systemDefaults: sd,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	mgmt_grpc.RegisterManagementServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return mgmtName
}

func (s *Server) MethodPrefix() string {
	return mgmt_grpc.ManagementService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return mgmt_grpc.ManagementService_AuthMethods
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
