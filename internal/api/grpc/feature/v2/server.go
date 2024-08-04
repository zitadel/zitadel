package feature

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

type Server struct {
	feature.UnimplementedFeatureServiceServer
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
	feature.RegisterFeatureServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return feature.FeatureService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return feature.FeatureService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return feature.FeatureService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return feature.RegisterFeatureServiceHandler
}
