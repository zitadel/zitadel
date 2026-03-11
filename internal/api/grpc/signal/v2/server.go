package signal

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	sig "github.com/zitadel/zitadel/internal/signals"
	signalpb "github.com/zitadel/zitadel/pkg/grpc/signal/v2"
	"github.com/zitadel/zitadel/pkg/grpc/signal/v2/signalconnect"
)

var _ signalconnect.SignalServiceHandler = (*Server)(nil)

// Server implements the SignalService connectRPC handler.
type Server struct {
	store *sig.DuckLakeStore
}

// CreateServer returns a new Server. Returns nil if the DuckLake store is unavailable.
func CreateServer(store *sig.DuckLakeStore) *Server {
	if store == nil {
		return nil
	}
	return &Server{store: store}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return signalconnect.NewSignalServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return signalpb.File_zitadel_signal_v2_signal_service_proto
}

func (s *Server) AppName() string {
	return signalpb.SignalService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return signalpb.SignalService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return signalpb.SignalService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return signalpb.RegisterSignalServiceHandler
}
