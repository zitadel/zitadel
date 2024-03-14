package schema

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	schema "github.com/zitadel/zitadel/pkg/grpc/user/schema/v3alpha"
)

var _ schema.UserSchemaServiceServer = (*Server)(nil)

type Server struct {
	schema.UnimplementedUserSchemaServiceServer
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
	schema.RegisterUserSchemaServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return schema.UserSchemaService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return schema.UserSchemaService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return schema.UserSchemaService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return schema.RegisterUserSchemaServiceHandler
}
