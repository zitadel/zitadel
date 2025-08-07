package action

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2/actionconnect"
)

var _ actionconnect.ActionServiceHandler = (*Server)(nil)

type Server struct {
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

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return actionconnect.NewActionServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return action.File_zitadel_action_v2_action_service_proto
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
