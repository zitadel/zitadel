package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// ReportQuotaUsage writes a slice of *quota.NotificationDueEvent directly to the eventstore
func (c *Commands) ReportQuotaUsage(ctx context.Context, dueNotifications []*quota.NotificationDueEvent) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	cmds := make([]eventstore.Command, 0, len(dueNotifications))
	for _, notification := range dueNotifications {
		ctxFilter, spanFilter := tracing.NewNamedSpan(ctx, "filterNotificationDueEvents")
		events, errFilter := c.eventstore.Filter(
			ctxFilter,
			eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				InstanceID(notification.Aggregate().InstanceID).
				AddQuery().
				AggregateTypes(quota.AggregateType).
				AggregateIDs(notification.Aggregate().ID).
				EventTypes(quota.NotificationDueEventType).
				EventData(map[string]interface{}{
					"id":          notification.ID,
					"periodStart": notification.PeriodStart,
					"threshold":   notification.Threshold,
				}).Builder(),
		)
		spanFilter.EndWithError(errFilter)
		if errFilter != nil {
			return errFilter
		}
		if len(events) > 0 {
			continue
		}
		cmds = append(cmds, notification)
	}
	if len(cmds) == 0 {
		return nil
	}
	ctxPush, spanPush := tracing.NewNamedSpan(ctx, "pushNotificationDueEvents")
	_, errPush := c.eventstore.Push(ctxPush, cmds...)
	spanPush.EndWithError(errPush)
	return errPush
}

func (c *Commands) UsageNotificationSent(ctx context.Context, dueEvent *quota.NotificationDueEvent) error {
	id, err := id_generator.Next()
	if err != nil {
		return err
	}

	_, err = c.eventstore.Push(
		ctx,
		quota.NewNotifiedEvent(ctx, id, dueEvent),
	)
	return err
}
