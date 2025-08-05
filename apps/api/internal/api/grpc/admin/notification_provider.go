package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/settings"
	"github.com/zitadel/zitadel/internal/domain"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetFileSystemNotificationProvider(ctx context.Context, req *admin_pb.GetFileSystemNotificationProviderRequest) (*admin_pb.GetFileSystemNotificationProviderResponse, error) {
	result, err := s.query.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeFile)
	if err != nil {
		return nil, err

	}
	return &admin_pb.GetFileSystemNotificationProviderResponse{
		Provider: settings.NotificationProviderToPb(result),
	}, nil
}

func (s *Server) GetLogNotificationProvider(ctx context.Context, req *admin_pb.GetLogNotificationProviderRequest) (*admin_pb.GetLogNotificationProviderResponse, error) {
	result, err := s.query.NotificationProviderByIDAndType(ctx, authz.GetInstance(ctx).InstanceID(), domain.NotificationProviderTypeLog)
	if err != nil {
		return nil, err

	}
	return &admin_pb.GetLogNotificationProviderResponse{
		Provider: settings.NotificationProviderToPb(result),
	}, nil
}
