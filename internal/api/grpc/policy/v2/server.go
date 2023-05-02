package policy

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	policy "github.com/zitadel/zitadel/pkg/grpc/policy/v2alpha"
)

var _ policy.PolicyServiceServer = (*Server)(nil)

type Server struct {
	policy.UnimplementedPolicyServiceServer
	command         *command.Commands
	query           *query.Queries
	assetsAPIDomain func(context.Context) string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	externalSecure bool,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		assetsAPIDomain: assets.AssetAPI(externalSecure),
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	policy.RegisterPolicyServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return policy.PolicyService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return policy.PolicyService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return policy.PolicyService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return policy.RegisterPolicyServiceHandler
}
