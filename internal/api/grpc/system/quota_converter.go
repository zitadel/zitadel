package system

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/quota"
)

type setQuotaRequest interface {
	GetUnit() quota.Unit
	GetFrom() *timestamppb.Timestamp
	GetResetInterval() *durationpb.Duration
	GetAmount() uint64
	GetLimit() bool
	GetNotifications() []*quota.Notification
}

func instanceQuotaPbToCommand(req setQuotaRequest) *command.SetQuota {
	return &command.SetQuota{
		Unit:          instanceQuotaUnitPbToCommand(req.GetUnit()),
		From:          req.GetFrom().AsTime(),
		ResetInterval: req.GetResetInterval().AsDuration(),
		Amount:        req.GetAmount(),
		Limit:         req.GetLimit(),
		Notifications: instanceQuotaNotificationsPbToCommand(req.GetNotifications()),
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
