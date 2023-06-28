package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/milestone"
)

// MilestonePushed writes a new milestone.PushedEvent with a new milestone.Aggregate to the eventstore
func (c *Commands) MilestonePushed(
	ctx context.Context,
	instanceID string,
	msType milestone.Type,
	endpoints []string,
	primaryDomain string,
) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, milestone.NewPushedEvent(ctx, milestone.NewAggregate(id, instanceID, instanceID), msType, endpoints, primaryDomain))
	return err
}
