package group

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	group "github.com/zitadel/zitadel/pkg/grpc/group/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/group/v2beta/groupconnect"
)

var _ groupconnect.GroupServiceHandler = (*Server)(nil)

type Server struct {
	systemDefaults systemdefaults.SystemDefaults
	command        *command.Commands
	query          *query.Queries

	checkPermission domain.PermissionCheck
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return groupconnect.NewGroupServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return group.File_zitadel_group_v2beta_group_service_proto
}

func (s *Server) AppName() string {
	return group.GroupService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return group.GroupService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return group.GroupService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return group.RegisterGroupServiceHandler
}
