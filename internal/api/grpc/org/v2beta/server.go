package org

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
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2beta/orgconnect"
)

var _ orgconnect.OrganizationServiceHandler = (*Server)(nil)

type Server struct {
	systemDefaults  systemdefaults.SystemDefaults
	command         *command.Commands
	query           *query.Queries
	checkPermission domain.PermissionCheck
}

type Config struct{}

func CreateServer(
	systemDefaults systemdefaults.SystemDefaults,
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
) *Server {
	return &Server{
		systemDefaults:  systemDefaults,
		command:         command,
		query:           query,
		checkPermission: checkPermission,
	}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return orgconnect.NewOrganizationServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return org.File_zitadel_org_v2beta_org_service_proto
}

func (s *Server) AppName() string {
	return org.OrganizationService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return org.OrganizationService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return org.OrganizationService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return org.RegisterOrganizationServiceHandler
}
