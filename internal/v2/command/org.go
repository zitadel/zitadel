package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) getOrg(ctx context.Context, orgID string) (*domain.Org, error) {
	writeModel, err := r.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if writeModel.State == domain.OrgStateActive {
		return nil, caos_errs.ThrowInternal(err, "COMMAND-4M9sf", "Errors.Org.NotFound")
	}
	return orgWriteModelToOrg(writeModel), nil
}

func (r *CommandSide) SetUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.User) error {
	orgAgg, userAgg, orgMemberAgg, err := r.setUpOrg(ctx, organisation, admin)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushAggregates(ctx, orgAgg, userAgg, orgMemberAgg)
	return err
}

func (r *CommandSide) setUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.User) (*org.Aggregate, *user.Aggregate, *org.Aggregate, error) {
	orgAgg, _, err := r.addOrg(ctx, organisation)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg, _, err := r.addHuman(ctx, orgAgg.ID(), admin.UserName, admin.Human)
	if err != nil {
		return nil, nil, nil, err
	}

	addedMember := NewOrgMemberWriteModel(orgAgg.ID(), userAgg.ID())
	orgMemberAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	err = r.addOrgMember(ctx, orgMemberAgg, addedMember, domain.NewMember(orgMemberAgg.ID(), userAgg.ID(), domain.OrgOwnerRole)) //TODO: correct?
	if err != nil {
		return nil, nil, nil, err
	}
	return orgAgg, userAgg, orgMemberAgg, nil
}

func (r *CommandSide) addOrg(ctx context.Context, organisation *domain.Org) (_ *org.Aggregate, _ *OrgWriteModel, err error) {
	if organisation == nil || !organisation.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
	}

	organisation.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, caos_errs.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}
	organisation.AddIAMDomain(r.iamDomain)
	addedOrg := NewOrgWriteModel(organisation.AggregateID)

	orgAgg := OrgAggregateFromWriteModel(&addedOrg.WriteModel)
	//TODO: uniqueness org name
	orgAgg.PushEvents(org.NewOrgAddedEvent(ctx, organisation.Name))
	for _, orgDomain := range organisation.Domains {
		if err := r.addOrgDomain(ctx, orgAgg, NewOrgDomainWriteModel(orgAgg.ID(), orgDomain.Domain), orgDomain); err != nil {
			return nil, nil, err
		}
	}
	return orgAgg, addedOrg, nil
}

func (r *CommandSide) getOrgWriteModelByID(ctx context.Context, orgID string) (*OrgWriteModel, error) {
	orgWriteModel := NewOrgWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModel, nil
}
