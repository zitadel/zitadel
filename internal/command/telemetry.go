package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/repository/milestone"
)

// ReportMilestoneReached writes each *milestone.ReachedEvent directly to the event store
func (c *Commands) ReportMilestoneReached(ctx context.Context, triggeringEvent eventstore.Event, customContext interface{}) error {
	aggregateId, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, milestone.NewReachedEvent(ctx, aggregateId, triggeringEvent, customContext))
	return err
}

// ReportMilestonePushed defers a milestone.PushedEvent for each *milestone.ReachedEvent and writes it directly to the event store.
func (c *Commands) ReportMilestonePushed(ctx context.Context, endpoints []string, reachedEvent *milestone.ReachedEvent) error {
	_, err := c.eventstore.Push(ctx, milestone.NewPushedEvent(ctx, reachedEvent, endpoints))
	return err
}
