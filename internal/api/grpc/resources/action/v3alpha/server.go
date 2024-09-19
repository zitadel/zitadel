package action

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/server"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
	action "github.com/zitadel/zitadel/v2/pkg/grpc/resources/action/v3alpha"
)

var _ action.ZITADELActionsServer = (*Server)(nil)

type Server struct {
	action.UnimplementedZITADELActionsServer
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
	action.RegisterZITADELActionsServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return action.ZITADELActions_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return action.ZITADELActions_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return action.ZITADELActions_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return action.RegisterZITADELActionsHandler
}

func checkActionsEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().Actions {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "ACTION-8o6pvqfjhs", "Errors.Action.NotEnabled")
}
