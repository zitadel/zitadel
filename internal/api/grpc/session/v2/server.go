package session

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2/sessionconnect"
)

var _ sessionconnect.SessionServiceHandler = (*Server)(nil)

type Server struct {
	command *command.Commands
	query   *query.Queries

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

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return sessionconnect.NewSessionServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return session.File_zitadel_session_v2_session_service_proto
}

func (s *Server) AppName() string {
	return session.SessionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return session.SessionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return session.SessionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return session.RegisterSessionServiceHandler
}

func (s *Server) GetFederatedLogoutRequest(ctx context.Context, req *connect.Request[session.GetFederatedLogoutRequestRequest]) (*connect.Response[session.GetFederatedLogoutRequestResponse], error) {
	logoutRequest, err := s.command.GetFederatedLogoutRequest(ctx, req.Msg.GetLogoutId())
	if err != nil {
		return nil, err
	}

	response := &session.GetFederatedLogoutRequestResponse{
		PostLogoutRedirectUri: logoutRequest.PostLogoutRedirectURI,
	}

	// Set the appropriate logout method based on binding type
	if logoutRequest.SAMLBindingType == "redirect" {
		response.LogoutMethod = &session.GetFederatedLogoutRequestResponse_Redirect{
			Redirect: &session.RedirectLogout{
				RedirectUri: logoutRequest.SAMLRedirectURL,
			},
		}
	} else if logoutRequest.SAMLBindingType == "post" {
		formData := map[string]string{
			"SAMLRequest": logoutRequest.SAMLRequest,
			"RelayState":  logoutRequest.SAMLRelayState,
		}
		response.LogoutMethod = &session.GetFederatedLogoutRequestResponse_Post{
			Post: &session.PostLogout{
				Url:      logoutRequest.SAMLPostURL,
				FormData: formData,
			},
		}
	}

	return connect.NewResponse(response), nil
}
