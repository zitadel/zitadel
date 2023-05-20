package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) SetSecurityPolicy(ctx context.Context, enabled bool, allowedOrigins []string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareSetSecurityPolicy(instanceAgg, enabled, allowedOrigins)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareSetSecurityPolicy(a *instance.Aggregate, enabled bool, allowedOrigins []string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getSecurityPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			cmd, err := writeModel.NewSetEvent(ctx, &a.Aggregate, enabled, allowedOrigins)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{cmd}, nil
		}, nil
	}
}

func (c *Commands) getSecurityPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (_ *InstanceSecurityPolicyWriteModel, err error) {
	writeModel := NewInstanceSecurityPolicyWriteModel(ctx)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
