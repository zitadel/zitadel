package execution

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

var _ execution.ExecutionServiceServer = (*Server)(nil)

type Server struct {
	execution.UnimplementedExecutionServiceServer
	command *command.Commands
	query   *query.Queries
}

type Config struct{}

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
	execution.RegisterExecutionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return execution.ExecutionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return execution.ExecutionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return execution.ExecutionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return execution.RegisterExecutionServiceHandler
}
