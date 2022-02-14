package settings

import (
	obj_pb "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
)

func NotificationProviderToPb(provider *query.DebugNotificationProvider) *settings_pb.DebugNotificationProvider {
	mapped := &settings_pb.DebugNotificationProvider{
		Compact: provider.Compact,
		Enabled: stateToEnabledAtt(provider.State),
		Details: obj_pb.ToViewDetailsPb(provider.Sequence, provider.CreationDate, provider.ChangeDate, provider.AggregateID),
	}
	return mapped
}

func stateToEnabledAtt(state domain.NotificationProviderState) bool {
	if state == domain.NotificationProviderStateEnabled {
		return true
	}
	return false
}
