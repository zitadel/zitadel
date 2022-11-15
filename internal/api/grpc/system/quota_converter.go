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
		Limitations:   instanceQuotaLimitationsPbToQuotaLimitations(req.Actions),
		Notifications: instanceQuotaNotificationsPbToQuotaNotifications(req.Notifications),
	}
}

func instanceQuotaLimitationsPbToQuotaLimitations(req *system_pb.AddQuotaRequest_Actions) *command.QuotaLimitations {
	return &command.QuotaLimitations{
		Block:       command.QuotaLimitationBlock{},
		CookieValue: req.CookieValue,
		RedirectURL: req.RedirectUrl,
	}
}

func instanceQuotaNotificationsPbToQuotaNotifications(req []*system_pb.AddQuotaRequest_Notification) command.QuotaNotifications {
	notifications := make([]*command.QuotaNotification, len(req))
	for idx := range req {
		item := req[idx]
		notifications[idx] = &command.QuotaNotification{
			Percent: item.Percent,
			Repeat:  item.Repeat,
			CallURL: item.CallUrl,
		}
	}
	return notifications
}
