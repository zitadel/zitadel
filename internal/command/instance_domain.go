package command

import (
	"context"

	"github.com/caos/logging"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddInstanceDomain(ctx context.Context, domain *domain.InstanceDomain) (*domain.InstanceDomain, error) {
	if !domain.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-R24hb", "Errors.Instance.InvalidDomain")
	}
	domainWriteModel := NewInstanceDomainWriteModel(domain.Domain)
	instanceAgg := IAMAggregateFromWriteModel(&domainWriteModel.WriteModel)
	events, err := c.addInstanceDomain(ctx, instanceAgg, domainWriteModel, domain, claimedUserIDs)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return domainWriteModelToOrgDomain(domainWriteModel), nil
}

func (c *Commands) RemoveInstanceDomain(ctx context.Context, instanceDomain *domain.InstanceDomain) (*domain.ObjectDetails, error) {
	if instanceDomain == nil || !instanceDomain.IsValid() || instanceDomain.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-SJsK3", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, instanceDomain.AggregateID, instanceDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Primary {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
	}
	instanceAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewDomainRemovedEvent(ctx, instanceAgg, instanceDomain.Domain, domainWriteModel.Verified))
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

	return org.NewDomainAddedEvent(ctx, instanceAgg, instanceDomain.Domain), nil
}

func (c *Commands) getInstanceDomainWriteModel(ctx context.Context, orgID, domain string) (*OrgDomainWriteModel, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgID, domain)
	err := c.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
