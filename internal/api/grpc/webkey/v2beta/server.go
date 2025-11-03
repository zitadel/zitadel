package webkey

import (
	"net/http"

	"connectrpc.com/connect"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta/webkeyconnect"
)

var _ webkeyconnect.WebKeyServiceHandler = (*Server)(nil)

type Server struct {
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

func (s *Server) AppName() string {
	return webkey.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return webkey.WebKeyService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return webkey.WebKeyService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return webkey.RegisterWebKeyServiceHandler
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return webkeyconnect.NewWebKeyServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return webkey.File_zitadel_webkey_v2beta_webkey_service_proto
}
