package settings

import (
	obj_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func NotificationProviderToPb(provider *query.DebugNotificationProvider) *settings_pb.DebugNotificationProvider {
	mapped := &settings_pb.DebugNotificationProvider{
		Compact: provider.Compact,
		Details: obj_pb.ToViewDetailsPb(provider.Sequence, provider.CreationDate, provider.ChangeDate, provider.AggregateID),
	}
	return mapped
}
