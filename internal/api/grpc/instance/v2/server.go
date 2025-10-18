package instance

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2/instanceconnect"
)

var _ instanceconnect.InstanceServiceHandler = (*Server)(nil)

type Server struct {
	command         *command.Commands
	query           *query.Queries
	systemDefaults  systemdefaults.SystemDefaults
	defaultInstance command.InstanceSetup
	externalDomain  string
	permissionCheck domain.PermissionCheck
}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	defaultInstance command.InstanceSetup,
	externalDomain string,
	check domain.PermissionCheck,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		defaultInstance: defaultInstance,
		externalDomain:  externalDomain,
		permissionCheck: check,
	}
}

func (s *Server) RegisterConnectServer(interceptors ...connect.Interceptor) (string, http.Handler) {
	return instanceconnect.NewInstanceServiceHandler(s, connect.WithInterceptors(interceptors...))
}

func (s *Server) FileDescriptor() protoreflect.FileDescriptor {
	return instance.File_zitadel_instance_v2_instance_service_proto
}

func (s *Server) AppName() string {
	return instance.InstanceService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return instance.InstanceService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return instance.InstanceService_AuthMethods
}

// checkPermission checks if either the system-wide or the instance-specific permission is granted.
func (s *Server) checkPermission(ctx context.Context, systemPermission, instancePermission string) error {
	// Let's first check the system permission since it's already resolved into the context.
	// If that succeeds, we don't need to resolve the instance permission.
	if err := s.permissionCheck(ctx, systemPermission, "", ""); err == nil {
		return nil
	}
	return s.permissionCheck(ctx, instancePermission, "", "")
}
