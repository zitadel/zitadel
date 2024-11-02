package webkey

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	webkey "github.com/zitadel/zitadel/pkg/grpc/resources/webkey/v3alpha"
)

type Server struct {
	webkey.UnimplementedZITADELWebKeysServer
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
	webkey.RegisterZITADELWebKeysServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return webkey.ZITADELWebKeys_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return webkey.ZITADELWebKeys_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return webkey.ZITADELWebKeys_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return webkey.RegisterZITADELWebKeysHandler
}
