package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

var (
	_ eventstore.Reducer    = (*RemoveOrg)(nil)
	_ eventstore.PushIntent = (*RemoveOrg)(nil)
)

type RemoveOrg struct {
	aggregate *eventstore.Aggregate
	commands  []eventstore.Command

	ID string
}

func NewRemoveOrg(id string) *RemoveOrg {
	return &RemoveOrg{
		ID: id,
	}
}

func (c *RemoveOrg) ToPushIntent(ctx context.Context) (eventstore.PushIntent, error) {
	c.aggregate = org.NewAggregate(ctx, c.ID)
	c.commands = append(c.commands, org.NewRemovedEvent(ctx))

	return c, nil
}

// Aggregate implements [eventstore.PushIntent].
func (c *RemoveOrg) Aggregate() *eventstore.Aggregate {
	return c.aggregate
}

// Commands implements [eventstore.PushIntent].
func (c *RemoveOrg) Commands() []eventstore.Command {
	return c.commands
}

// CurrentSequence implements [eventstore.PushIntent].
func (*RemoveOrg) CurrentSequence() eventstore.CurrentSequence {
	// TODO: implement after filter works
	return eventstore.SequenceAtLeast(1)
}

// Reduce implements [eventstore.Reducer].
func (*RemoveOrg) Reduce(events ...eventstore.Event) error {
	panic("unimplemented")
}
