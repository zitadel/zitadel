package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddOrgIAMPolicy(ctx context.Context, resourceOwner string, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	addedPolicy := NewORGOrgIAMPolicyWriteModel(resourceOwner)
	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.PolicyOrgIAMWriteModel.WriteModel)
	event, err := r.addOrgIAMPolicy(ctx, orgAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrgIAMPolicy(addedPolicy), nil
}

func (r *CommandSide) addOrgIAMPolicy(ctx context.Context, orgAgg *eventstore.Aggregate, addedPolicy *ORGOrgIAMPolicyWriteModel, policy *domain.OrgIAMPolicy) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-1M8ds", "Errors.Org.OrgIAMPolicy.AlreadyExists")
	}
	return org.NewOrgIAMPolicyAddedEvent(ctx, orgAgg, policy.UserLoginMustBeDomain), nil
}

func (r *CommandSide) ChangeOrgIAMPolicy(ctx context.Context, resourceOwner string, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	existingPolicy, err := r.orgIAMPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-2N9sd", "Errors.Org.OrgIAMPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-3M9ds", "Errors.Org.LabelPolicy.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrgIAMPolicy(existingPolicy), nil
}

func (r *CommandSide) RemoveOrgIAMPolicy(ctx context.Context, orgID string) error {
	existingPolicy, err := r.orgIAMPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-Dvsh3", "Errors.Org.OrgIAMPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, org.NewOrgIAMPolicyRemovedEvent(ctx, orgAgg))
	return err
}

func (r *CommandSide) getOrgIAMPolicy(ctx context.Context, orgID string) (*domain.OrgIAMPolicy, error) {
	policy, err := r.orgIAMPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return orgWriteModelToOrgIAMPolicy(policy), nil
	}
	return r.getDefaultOrgIAMPolicy(ctx)
}

func (r *CommandSide) orgIAMPolicyWriteModelByID(ctx context.Context, orgID string) (policy *ORGOrgIAMPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewORGOrgIAMPolicyWriteModel(orgID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
