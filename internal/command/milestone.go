package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

// MilestonePushed writes a new milestone.PushedEvent with a new milestone.Aggregate to the eventstore
func (c *Commands) MilestonePushed(
	ctx context.Context,
	msType milestone.Type,
	endpoints []string,
	primaryDomain string,
) error {
	id, err := id_generator.Next()
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, milestone.NewPushedEvent(ctx, milestone.NewAggregate(ctx, id), msType, endpoints, primaryDomain, c.externalDomain))
	return err
}
