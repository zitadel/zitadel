package template

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/pkg/grpc/settings/template/v2"
)

var _ template.TemplateServiceServer = (*Server)(nil)

type Server struct {
	template.UnimplementedTemplateServiceServer
}

func CreateServer() *Server {
	return &Server{}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	template.RegisterTemplateServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return template.TemplateService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return template.TemplateService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping { return template.TemplateService_AuthMethods }

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return template.RegisterTemplateServiceHandler
}
