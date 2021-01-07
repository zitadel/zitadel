package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.OrgDomain, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgDomain.AggregateID, orgDomain.Domain)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	err := r.addOrgDomain(ctx, orgAgg, domainWriteModel, orgDomain)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
	if err != nil {
		return nil, err
	}
	return orgDomainWriteModelToOrgDomain(domainWriteModel), nil
}

func (r *CommandSide) addOrgDomain(ctx context.Context, orgAgg *org.Aggregate, addedDomain *OrgDomainWriteModel, orgDomain *domain.OrgDomain) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedDomain)
	if err != nil {
		return err
	}
	if addedDomain.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "COMMA-Bd2jj", "Errors.Org.Domain.AlreadyExists")
	}
	orgAgg.PushEvents(org.NewDomainAddedEvent(ctx, orgDomain.Domain))
	if orgDomain.Verified {
		//TODO: uniqueness verified domain
		//TODO: users with verified domain -> domain claimed
		orgAgg.PushEvents(org.NewDomainVerifiedEvent(ctx, orgDomain.Domain))
	}
	if orgDomain.Primary {
		orgAgg.PushEvents(org.NewDomainPrimarySetEvent(ctx, orgDomain.Domain))
	}
	return nil
}
