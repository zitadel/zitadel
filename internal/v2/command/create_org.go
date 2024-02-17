package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

var (
	_ eventstore.Reducer    = (*CreateOrg)(nil)
	_ eventstore.PushIntent = (*CreateOrg)(nil)
)

type CreateOrg struct {
	aggregate *eventstore.Aggregate
	commands  []eventstore.Command

	Name   string
	Domain string
}

func NewCreateOrg(name string) *CreateOrg {
	return &CreateOrg{
		Name: name,
	}
}

func (c *CreateOrg) ToPushIntent(ctx context.Context) (eventstore.PushIntent, error) {
	orgID, err := id.SonyFlakeGenerator().Next()
	if err != nil {
		return nil, err
	}
	c.aggregate = org.NewAggregate(ctx, orgID)

	c.Domain, err = domain.NewIAMDomainName(c.Name, authz.GetInstance(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}

	for _, err := range []error{
		c.appendCommand(org.NewAddedEvent(ctx, c.Name)),
		c.appendCommand(org.NewDomainAddedEvent(ctx, c.Domain)),
		c.appendCommand(org.NewDomainVerifiedEvent(ctx, c.Domain)),
		c.appendCommand(org.NewSetDomainPrimaryEvent(ctx, c.Domain)),
	} {
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *CreateOrg) appendCommand(command eventstore.Command, err error) error {
	if err != nil {
		return err
	}
	c.commands = append(c.commands, command)
	return nil
}

// Aggregate implements [eventstore.PushIntent].
func (c *CreateOrg) Aggregate() *eventstore.Aggregate {
	return c.aggregate
}

// Commands implements [eventstore.PushIntent].
func (c *CreateOrg) Commands() []eventstore.Command {
	return c.commands
}

// CurrentSequence implements [eventstore.PushIntent].
func (*CreateOrg) CurrentSequence() eventstore.CurrentSequence {
	return eventstore.SequenceMatches(0)
}

// Reduce implements [eventstore.Reducer].
func (*CreateOrg) Reduce(events ...eventstore.Event) error {
	panic("unimplemented")
}
