package settings

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2/settingsconnect"
)

var _ settingsconnect.SettingsServiceHandler = (*Server)(nil)

type Server struct {
	systemDefaults systemdefaults.SystemDefaults
	command        *command.Commands
	query          *query.Queries

	checkPermission domain.PermissionCheck
	assetsAPIDomain func(context.Context) string
	// smtpRules are the operator rules restricting email providers
	smtpRules smtp.Rules
}

func CreateServer(
	systemDefaults systemdefaults.SystemDefaults,
	command *command.Commands,
	query *query.Queries,
	checkPermission domain.PermissionCheck,
	smtpRules smtp.Rules,
) *Server {
	return &Server{
		systemDefaults:  systemDefaults,
		command:         command,
		query:           query,
		checkPermission: checkPermission,
		assetsAPIDomain: assets.AssetAPI(),
		smtpRules:       smtpRules,
	}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return settingsconnect.NewSettingsServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return settings.File_zitadel_settings_v2_settings_service_proto
}

func (s *Server) AppName() string {
	return settings.SettingsService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return settings.SettingsService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return settings.SettingsService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return settings.RegisterSettingsServiceHandler
}
