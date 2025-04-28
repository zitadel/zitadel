package instance

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

var _ instance.InstanceServiceServer = (*Server)(nil)

type Server struct {
	instance.UnimplementedInstanceServiceServer
	command         *command.Commands
	query           *query.Queries
	checkPermission domain.PermissionCheck
	systemDefaults  systemdefaults.SystemDefaults
	defaultInstance command.InstanceSetup
	externalDomain  string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
	systemDefaults systemdefaults.SystemDefaults,
	defaultInstance command.InstanceSetup,
	externalDomain string,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		checkPermission: checkPermission,
		systemDefaults:  systemDefaults,
		defaultInstance: defaultInstance,
		externalDomain:  externalDomain,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	instance.RegisterInstanceServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return instance.InstanceService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return instance.InstanceService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return instance.InstanceService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return instance.RegisterInstanceServiceHandler
}
