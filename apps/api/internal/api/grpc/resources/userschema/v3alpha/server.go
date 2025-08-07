package userschema

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

var _ schema.ZITADELUserSchemasServer = (*Server)(nil)

type Server struct {
	schema.UnimplementedZITADELUserSchemasServer
	systemDefaults systemdefaults.SystemDefaults
	command        *command.Commands
	query          *query.Queries
}

type Config struct{}

func CreateServer(
	systemDefaults systemdefaults.SystemDefaults,
	command *command.Commands,
	query *query.Queries,
) *Server {
	return &Server{
		systemDefaults: systemDefaults,
		command:        command,
		query:          query,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	schema.RegisterZITADELUserSchemasServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return schema.ZITADELUserSchemas_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return schema.ZITADELUserSchemas_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return schema.ZITADELUserSchemas_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return schema.RegisterZITADELUserSchemasHandler
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "SCHEMA-SFjk3", "Errors.UserSchema.NotEnabled")
}
