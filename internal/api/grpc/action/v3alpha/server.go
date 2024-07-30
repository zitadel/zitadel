package action

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

var _ action.ActionServiceServer = (*Server)(nil)

type Server struct {
	action.UnimplementedActionServiceServer
	command             *command.Commands
	query               *query.Queries
	ListActionFunctions func() []string
	ListGRPCMethods     func() []string
	ListGRPCServices    func() []string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	listActionFunctions func() []string,
	listGRPCMethods func() []string,
	listGRPCServices func() []string,
) *Server {
	return &Server{
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

func checkExecutionEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().Actions {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "ACTION-8o6pvqfjhs", "Errors.Action.NotEnabled")
}
