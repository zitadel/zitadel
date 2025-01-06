package saml

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
)

var _ saml_pb.SAMLServiceServer = (*Server)(nil)

type Server struct {
	saml_pb.UnimplementedSAMLServiceServer
	command *command.Commands
	query   *query.Queries

	idp            *saml.Provider
	externalSecure bool
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	idp *saml.Provider,
	externalSecure bool,
) *Server {
	return &Server{
		command:        command,
		query:          query,
		idp:            idp,
		externalSecure: externalSecure,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	saml_pb.RegisterSAMLServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return saml_pb.SAMLService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return saml_pb.SAMLService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return saml_pb.SAMLService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return saml_pb.RegisterSAMLServiceHandler
}
