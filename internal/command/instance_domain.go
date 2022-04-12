package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func (c *Commands) AddInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.addInstanceDomain(instanceAgg, instanceDomain, false)
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
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) RemoveInstanceDomain(ctx context.Context, instanceDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.removeInstanceDomain(instanceAgg, instanceDomain)
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
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) addInstanceDomain(a *instance.Aggregate, instanceDomain string, generated bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-28nlD", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainWriteModel, err := c.getInstanceDomainWriteModel(ctx, instanceDomain)
			if err != nil {
				return nil, err
			}
			if domainWriteModel.State == domain.InstanceDomainStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-i2nl", "Errors.Instance.Domain.AlreadyExists")
			}
			return []eventstore.Command{instance.NewDomainAddedEvent(ctx, &a.Aggregate, instanceDomain, generated)}, nil
		}, nil
	}
}

func (c *Commands) removeInstanceDomain(a *instance.Aggregate, instanceDomain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if instanceDomain = strings.TrimSpace(instanceDomain); instanceDomain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-39nls", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainWriteModel, err := c.getInstanceDomainWriteModel(ctx, instanceDomain)
			if err != nil {
				return nil, err
			}
			if domainWriteModel.State != domain.InstanceDomainStateActive {
				return nil, errors.ThrowNotFound(nil, "INSTANCE-8ls9f", "Errors.Instance.Domain.NotFound")
			}
			if domainWriteModel.Generated {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTANCE-9hn3n", "Errors.Instance.Domain.GeneratedNotRemovable")
			}
			return []eventstore.Command{instance.NewDomainRemovedEvent(ctx, &a.Aggregate, instanceDomain)}, nil
		}, nil
	}
}

func (c *Commands) getInstanceDomainWriteModel(ctx context.Context, domain string) (*InstanceDomainWriteModel, error) {
	domainWriteModel := NewInstanceDomainWriteModel(ctx, domain)
	err := c.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
