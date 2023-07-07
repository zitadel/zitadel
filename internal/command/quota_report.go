package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

// ReportQuotaUsage writes a slice of *quota.NotificationDueEvent directly to the eventstore
func (c *Commands) ReportQuotaUsage(ctx context.Context, dueNotifications []*quota.NotificationDueEvent) error {
	cmds := make([]eventstore.Command, len(dueNotifications))
	for idx, notification := range dueNotifications {
		cmds[idx] = notification
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
