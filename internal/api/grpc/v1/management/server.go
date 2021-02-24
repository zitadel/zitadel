package management

import (
	"github.com/caos/zitadel/internal/v2/command"
	"github.com/caos/zitadel/internal/v2/query"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/management/repository"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing"
	"github.com/caos/zitadel/pkg/grpc/management"
)

const (
	mgmtName = "Management-API"
)

var _ management.ManagementServiceServer = (*Server)(nil)

type Server struct {
	command        *command.CommandSide
	query          *query.QuerySide
	project        repository.ProjectRepository
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

func CreateServer(command *command.CommandSide, query *query.QuerySide, repo repository.Repository, sd systemdefaults.SystemDefaults) *Server {
	return &Server{
		command:        command,
		query:          query,
		project:        repo,
		org:            repo,
		user:           repo,
		usergrant:      repo,
		iam:            repo,
		systemDefaults: sd,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	management.RegisterManagementServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return mgmtName
}

func (s *Server) MethodPrefix() string {
	return management.ManagementService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return management.ManagementService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return management.RegisterManagementServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/management/v1"
}
