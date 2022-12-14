package system

import (
	"time"

	"github.com/zitadel/zitadel/internal/command"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func instanceQuotaPbToQuota(req *system_pb.AddQuotaRequest) *command.Quota {
	return &command.Quota{
		Unit:          command.QuotaUnit(req.Unit),
		From:          req.From.AsTime().Format(time.RFC3339),
		Interval:      req.Interval.AsDuration(),
		Amount:        req.Amount,
		Limit:         req.Limit,
		Notifications: instanceQuotaNotificationsPbToQuotaNotifications(req.Notifications),
	}
}

func instanceQuotaNotificationsPbToQuotaNotifications(req []*system_pb.AddQuotaRequest_Notification) command.QuotaNotifications {
	notifications := make([]*command.QuotaNotification, len(req))
	for idx := range req {
		item := req[idx]
		notifications[idx] = &command.QuotaNotification{
			Percent: uint64(item.Percent),
			Repeat:  item.Repeat,
			CallURL: item.CallUrl,
		}
	}
	return notifications
}
