package action

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
)

var _ action.ActionServiceServer = (*Server)(nil)

type Server struct {
	action.UnimplementedActionServiceServer
	systemDefaults      systemdefaults.SystemDefaults
	command             *command.Commands
	query               *query.Queries
	ListActionFunctions func() []string
	ListGRPCMethods     func() []string
	ListGRPCServices    func() []string
}

type Config struct{}

func CreateServer(
	systemDefaults systemdefaults.SystemDefaults,
	command *command.Commands,
	query *query.Queries,
	listActionFunctions func() []string,
	listGRPCMethods func() []string,
	listGRPCServices func() []string,
) *Server {
	return &Server{
		systemDefaults:      systemDefaults,
		command:             command,
		query:               query,
		ListActionFunctions: listActionFunctions,
		ListGRPCMethods:     listGRPCMethods,
		ListGRPCServices:    listGRPCServices,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	action.RegisterActionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return action.ActionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return action.ActionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return action.ActionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return action.RegisterActionServiceHandler
}
