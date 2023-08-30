package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

// ReportQuotaUsage writes a slice of *quota.NotificationDueEvent directly to the eventstore
func (c *Commands) ReportQuotaUsage(ctx context.Context, dueNotifications []*quota.NotificationDueEvent) error {
	cmds := make([]eventstore.Command, 0)
	for _, notification := range dueNotifications {
		// TODO: doesnt' work, never returns any events
		events, err := c.eventstore.Filter(
			ctx,
			eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				InstanceID(notification.Aggregate().InstanceID).
				AddQuery().
				AggregateTypes(quota.AggregateType).
				AggregateIDs(notification.Aggregate().ID).
				EventTypes(quota.NotificationDueEventType).
				EventData(map[string]interface{}{
					"id":          notification.ID,
					"periodStart": notification.PeriodStart.Format("2023-08-30T02:00:00+02:00"),
					"threshold":   notification.Threshold,
				}).Builder(),
		)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			continue
		}
		cmds = append(cmds, notification)
	}
	if len(cmds) == 0 {
		return nil
	}
	_, err := c.eventstore.Push(ctx, cmds...)
	return err
}

func (c *Commands) UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	_, err = c.eventstore.Push(
		ctx,
		quota.NewNotifiedEvent(ctx, id, dueEvent),
	)
	return err
}
