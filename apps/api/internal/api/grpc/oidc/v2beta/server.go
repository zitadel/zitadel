package oidc

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta/oidcconnect"
)

var _ oidcconnect.OIDCServiceHandler = (*Server)(nil)

type Server struct {
	command *command.Commands
	query   *query.Queries

	op             *oidc.Server
	externalSecure bool
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	op *oidc.Server,
	externalSecure bool,
) *Server {
	return &Server{
		command:        command,
		query:          query,
		op:             op,
		externalSecure: externalSecure,
	}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return oidcconnect.NewOIDCServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return oidc_pb.File_zitadel_oidc_v2beta_oidc_service_proto
}

func (s *Server) AppName() string {
	return oidc_pb.OIDCService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return oidc_pb.OIDCService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return oidc_pb.OIDCService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return oidc_pb.RegisterOIDCServiceHandler
}
