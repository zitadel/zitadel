package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) getOrg(ctx context.Context, orgID string) (*domain.Org, error) {
	writeModel, err := r.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if writeModel.State == domain.OrgStateUnspecified || writeModel.State == domain.OrgStateRemoved {
		return nil, caos_errs.ThrowInternal(err, "COMMAND-4M9sf", "Errors.Org.NotFound")
	}
	return orgWriteModelToOrg(writeModel), nil
}

func (r *CommandSide) checkOrgExists(ctx context.Context, orgID string) error {
	orgWriteModel, err := r.getOrgWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.Org.NotFound")
	}
	return nil
}

func (r *CommandSide) SetUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.Human) error {
	_, _, _, events, err := r.setUpOrg(ctx, organisation, admin)
	if err != nil {
		return err
	}
	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) AddOrg(ctx context.Context, name, userID, resourceOwner string) (*domain.Org, error) {
	orgAgg, addedOrg, events, err := r.addOrg(ctx, &domain.Org{Name: name})
	if err != nil {
		return nil, err
	}

	err = r.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedMember := NewOrgMemberWriteModel(addedOrg.AggregateID, userID)
	orgMemberEvent, err := r.addOrgMember(ctx, orgAgg, addedMember, domain.NewMember(orgAgg.ID, userID, domain.RoleOrgOwner))
	if err != nil {
		return nil, err
	}
	events = append(events, orgMemberEvent)
	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedOrg, pushedEvents...)
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
	_, err = r.eventstore.PushEvents(ctx, org.NewOrgDeactivatedEvent(ctx, orgAgg))
	return err
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
	_, err = r.eventstore.PushEvents(ctx, org.NewOrgReactivatedEvent(ctx, orgAgg))
	return err
}

func (r *CommandSide) setUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.Human) (orgAgg *eventstore.Aggregate, human *HumanWriteModel, orgMember *OrgMemberWriteModel, events []eventstore.EventPusher, err error) {
	orgAgg, _, addOrgEvents, err := r.addOrg(ctx, organisation)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userEvents, human, err := r.addHuman(ctx, orgAgg.ID, admin)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	addOrgEvents = append(addOrgEvents, userEvents...)

	addedMember := NewOrgMemberWriteModel(orgAgg.ID, human.AggregateID)
	orgMemberAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	orgMemberEvent, err := r.addOrgMember(ctx, orgMemberAgg, addedMember, domain.NewMember(orgMemberAgg.ID, human.AggregateID, domain.RoleOrgOwner))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	addOrgEvents = append(addOrgEvents, orgMemberEvent)
	return orgAgg, human, addedMember, addOrgEvents, nil
}

func (r *CommandSide) addOrg(ctx context.Context, organisation *domain.Org, claimedUserIDs ...string) (_ *eventstore.Aggregate, _ *OrgWriteModel, _ []eventstore.EventPusher, err error) {
	if organisation == nil || !organisation.IsValid() {
		return nil, nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
	}

	organisation.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, caos_errs.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}
	organisation.AddIAMDomain(r.iamDomain)
	addedOrg := NewOrgWriteModel(organisation.AggregateID)

	orgAgg := OrgAggregateFromWriteModel(&addedOrg.WriteModel)
	events := []eventstore.EventPusher{
		org.NewOrgAddedEvent(ctx, orgAgg, organisation.Name),
	}
	for _, orgDomain := range organisation.Domains {
		orgDomainEvents, err := r.addOrgDomain(ctx, orgAgg, NewOrgDomainWriteModel(orgAgg.ID, orgDomain.Domain), orgDomain, claimedUserIDs...)
		if err != nil {
			return nil, nil, nil, err
		} else {
			events = append(events, orgDomainEvents...)
		}
	}
	return orgAgg, addedOrg, events, nil
}

func (r *CommandSide) getOrgWriteModelByID(ctx context.Context, orgID string) (*OrgWriteModel, error) {
	orgWriteModel := NewOrgWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModel, nil
}
