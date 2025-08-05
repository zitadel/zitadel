package debug_events

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	debug_events "github.com/zitadel/zitadel/pkg/grpc/resources/debug_events/v3alpha"
)

type Server struct {
	debug_events.UnimplementedZITADELDebugEventsServer
	command *command.Commands
	query   *query.Queries
}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
) *Server {
	return &Server{
		command: command,
		query:   query,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	debug_events.RegisterZITADELDebugEventsServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return debug_events.ZITADELDebugEvents_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return debug_events.ZITADELDebugEvents_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return debug_events.ZITADELDebugEvents_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return debug_events.RegisterZITADELDebugEventsHandler
}
