package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/milestone"
)

// MilestonePushed writes a new event with a new milestone.Aggregate to the eventstore
func (c *Commands) MilestonePushed(
	ctx context.Context,
	instanceID string,
	eventType milestone.PushedEventType,
	endpoints []string,
	primaryDomain string,
) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	pushedEvent, err := milestone.NewPushedEventByType(ctx, eventType, milestone.NewAggregate(id, instanceID, instanceID), endpoints, primaryDomain)
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, pushedEvent)
	return err
}
