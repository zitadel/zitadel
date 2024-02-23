package org

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

var _ org.OrganizationServiceServer = (*Server)(nil)

type Server struct {
	org.UnimplementedOrganizationServiceServer
	command         *command.Commands
	query           *query.Queries
	checkPermission domain.PermissionCheck

	es *eventstore.EventStore
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
	es *eventstore.EventStore,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		checkPermission: checkPermission,
		es:              es,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	org.RegisterOrganizationServiceServer(grpcServer, s)
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
