package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/milestone"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// ReportTelemetryUsage writes one or many *telemetry.PushDueEvent directly to the eventstore
func (c *Commands) ReportTelemetryUsage(ctx context.Context, dueEvent ...*milestone.ReachedEvent) error {
	cmds := make([]eventstore.Command, len(dueEvent))
	for idx, notification := range dueEvent {
		cmds[idx] = notification
	}
	_, err := c.eventstore.Push(ctx, cmds...)
	return err
}

func (c *Commands) TelemetryPushed(ctx context.Context, dueEvent *milestone.ReachedEvent, endpoints []string) error {
	_, err := c.eventstore.Push(
		ctx,
		milestone.NewPushedEvent(ctx, dueEvent, endpoints),
	)
	return err
}
