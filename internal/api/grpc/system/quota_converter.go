package system

import (
	"github.com/zitadel/zitadel/internal/command"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func instanceQuotaPbToQuota(req *system_pb.AddQuotaRequest) *command.AddQuota {
	return &command.AddQuota{
		Unit:          instanceQuotaUnitPbToQuotaUnit(req.Unit),
		From:          req.From.AsTime(),
		Interval:      req.Interval.AsDuration(),
		Amount:        req.Amount,
		Limit:         req.Limit,
		Notifications: instanceQuotaNotificationsPbToQuotaNotifications(req.Notifications),
	}
}

func instanceQuotaUnitPbToQuotaUnit(unit system_pb.Unit) command.QuotaUnit {
	switch unit {
	case system_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED:
		return command.QuotaRequestsAllAuthenticated
	case system_pb.Unit_UNIT_ACTIONS_ALL_RUN_SECONDS:
		return command.QuotaActionsAllRunsSeconds
	case system_pb.Unit_UNIT_UNIMPLEMENTED:
		fallthrough
	default:
		return command.QuotaUnit(unit.String())
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
