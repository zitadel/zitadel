package command

import (
	"context"

	"github.com/caos/zitadel/internal/repository/iam"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
)

func (c *Commands) AddInstanceDomain(ctx context.Context, domain *domain.InstanceDomain) (*domain.InstanceDomain, error) {
	if !domain.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-R24hb", "Errors.Instance.Domain.Invalid")
	}
	domainWriteModel := NewInstanceDomainWriteModel(domain.Domain)
	instanceAgg := IAMAggregateFromWriteModel(&domainWriteModel.WriteModel)
	event, err := c.addInstanceDomain(ctx, instanceAgg, domainWriteModel, domain)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return instanceDomainWriteModelToInstanceDomain(domainWriteModel), nil
}

func (c *Commands) RemoveInstanceDomain(ctx context.Context, instanceDomain *domain.InstanceDomain) (*domain.ObjectDetails, error) {
	if instanceDomain == nil || !instanceDomain.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-SJsK3", "Errors.Instance.Domain.Invalid")
	}
	domainWriteModel, err := c.getInstanceDomainWriteModel(ctx, instanceDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.InstanceDomainStateActive {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-8ls9f", "Errors.Instance.Domain.NotFound")
	}
	if domainWriteModel.Generated {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-9hn3n", "Errors.Instance.Domain.GeneratedNotRemovable")
	}
	instanceAgg := IAMAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewDomainRemovedEvent(ctx, instanceAgg, instanceDomain.Domain))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
}

func (c *Commands) addInstanceDomain(ctx context.Context, instanceAgg *eventstore.Aggregate, addedDomain *InstanceDomainWriteModel, instanceDomain *domain.InstanceDomain) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedDomain)
	if err != nil {
		return nil, err
	}
	if addedDomain.State == domain.InstanceDomainStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMA-nfske", "Errors.Instance.Domain.AlreadyExists")
	}

	return iam.NewDomainAddedEvent(ctx, instanceAgg, instanceDomain.Domain, instanceDomain.Generated), nil
}

func (c *Commands) getInstanceDomainWriteModel(ctx context.Context, domain string) (*InstanceDomainWriteModel, error) {
	domainWriteModel := NewInstanceDomainWriteModel(domain)
	err := c.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
