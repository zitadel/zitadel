package system

import (
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/quota"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func instanceQuotaPbToCommand(req *system.AddQuotaRequest) *command.AddQuota {
	return &command.AddQuota{
		Unit:          instanceQuotaUnitPbToCommand(req.Unit),
		From:          req.From.AsTime(),
		ResetInterval: req.ResetInterval.AsDuration(),
		Amount:        req.Amount,
		Limit:         req.Limit,
		Notifications: instanceQuotaNotificationsPbToCommand(req.Notifications),
	}
}

func instanceQuotaUnitPbToCommand(unit quota.Unit) command.QuotaUnit {
	switch unit {
	case quota.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED:
		return command.QuotaRequestsAllAuthenticated
	case quota.Unit_UNIT_ACTIONS_ALL_RUN_SECONDS:
		return command.QuotaActionsAllRunsSeconds
	case quota.Unit_UNIT_UNIMPLEMENTED:
		fallthrough
	default:
		return command.QuotaUnit(unit.String())
	}
}

func instanceQuotaNotificationsPbToCommand(req []*quota.Notification) command.QuotaNotifications {
	notifications := make([]*command.QuotaNotification, len(req))
	for idx, item := range req {
		notifications[idx] = &command.QuotaNotification{
			Percent: uint16(item.Percent),
			Repeat:  item.Repeat,
			CallURL: item.CallUrl,
		}
	}
	return notifications
}
