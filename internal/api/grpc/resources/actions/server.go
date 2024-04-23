package actions

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/pkg/grpc/resouces/action/v2"
)

var _ action.ActionServiceServer = (*Server)(nil)

type Server struct {
	action.UnimplementedActionServiceServer
}

func CreateServer() *Server {
	return &Server{}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	action.RegisterActionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return action.ActionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return action.ActionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping { return action.ActionService_AuthMethods }

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return action.RegisterActionServiceHandler
}
