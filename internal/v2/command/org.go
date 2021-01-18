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

func (r *CommandSide) SetUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.Human) error {
	orgAgg, userAgg, orgMemberAgg, err := r.setUpOrg(ctx, organisation, admin)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushAggregates(ctx, orgAgg, userAgg, orgMemberAgg)
	return err
}

func (r *CommandSide) AddOrg(ctx context.Context, name, userID, resourceOwner string) (*domain.Org, error) {
	orgAgg, addedOrg, err := r.addOrg(ctx, &domain.Org{Name: name})
	if err != nil {
		return nil, err
	}

	active, err := r.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, caos_errs.ThrowPreconditionFailed(err, "ORG-HBR2z", "Errors.User.NotFound")
	}
	addedMember := NewOrgMemberWriteModel(orgAgg.ID(), userID)
	err = r.addOrgMember(ctx, orgAgg, addedMember, domain.NewMember(orgAgg.ID(), userID, domain.RoleOrgOwner))
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedOrg, orgAgg)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrg(addedOrg), nil
}

func (r *CommandSide) DeactivateOrg(ctx context.Context, orgID string) error {
	orgWriteModel, err := r.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-oL9nT", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateInactive {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-Dbs2g", "Errors.Org.AlreadyDeactivated")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewOrgDeactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, orgWriteModel, orgAgg)
}

func (r *CommandSide) ReactivateOrg(ctx context.Context, orgID string) error {
	orgWriteModel, err := r.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-Dgf3g", "Errors.Org.NotFound")
	}
	if orgWriteModel.State == domain.OrgStateActive {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-bfnrh", "Errors.Org.AlreadyActive")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewOrgReactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, orgWriteModel, orgAgg)
}

func (r *CommandSide) setUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.Human) (*org.Aggregate, *user.Aggregate, *org.Aggregate, error) {
	orgAgg, _, err := r.addOrg(ctx, organisation)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg, _, err := r.addHuman(ctx, orgAgg.ID(), admin)
	if err != nil {
		return nil, nil, nil, err
	}

	addedMember := NewOrgMemberWriteModel(orgAgg.ID(), userAgg.ID())
	orgMemberAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	err = r.addOrgMember(ctx, orgMemberAgg, addedMember, domain.NewMember(orgMemberAgg.ID(), userAgg.ID(), domain.RoleOrgOwner))
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
