package signal

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/signals"
	signal "github.com/zitadel/zitadel/pkg/grpc/signal/v2"
	"github.com/zitadel/zitadel/pkg/grpc/signal/v2/signalconnect"
)

var _ signalconnect.SignalServiceHandler = (*Server)(nil)

type Server struct {
	reader signals.SignalReader
}

func CreateServer(reader signals.SignalReader) *Server {
	if reader == nil {
		return nil
	}
	return &Server{
		reader: reader,
	}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return signalconnect.NewSignalServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return signal.File_zitadel_signal_v2_signal_service_proto
}

func (s *Server) AppName() string {
	return signal.SignalService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return signal.SignalService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return signal.SignalService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return signal.RegisterSignalServiceHandler
}
