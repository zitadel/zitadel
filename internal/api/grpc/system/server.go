package system

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/admin/repository"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/system"
)

const (
	systemAPI = "System-API"
)

var _ system.SystemServiceServer = (*Server)(nil)

type Server struct {
	system.UnimplementedSystemServiceServer
	command         *command.Commands
	query           *query.Queries
	administrator   repository.AdministratorRepository
	DefaultInstance command.InstanceSetup
	ExternalSecure  bool
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	defaultInstance command.InstanceSetup,
	externalSecure bool) *Server {
	return &Server{
		command:         command,
		query:           query,
		administrator:   repo,
		DefaultInstance: defaultInstance,
		ExternalSecure:  externalSecure,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	system.RegisterSystemServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return systemAPI
}

func (s *Server) MethodPrefix() string {
	return system.SystemService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return system.SystemService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return system.RegisterSystemServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/system/v1"
}
