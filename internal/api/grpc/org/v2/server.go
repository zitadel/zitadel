package org

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

var _ org.OrganisationServiceServer = (*Server)(nil)

type Server struct {
	org.UnimplementedOrganisationServiceServer
	command         *command.Commands
	query           *query.Queries
	checkPermission domain.PermissionCheck
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		checkPermission: checkPermission,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	org.RegisterOrganisationServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return org.OrganisationService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return org.OrganisationService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return org.OrganisationService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return org.RegisterOrganisationServiceHandler
}
