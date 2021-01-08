package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	addedPolicy := NewORGOrgIAMPolicyWriteModel(policy.AggregateID)
	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.PolicyOrgIAMWriteModel.WriteModel)
	err := r.addOrgIAMPolicy(ctx, orgAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-5M0ds", "Errors.Org.OrgIAMPolicy.AlreadyExists")
	}
	orgAgg.PushEvents(org.NewOrgIAMPolicyAddedEvent(ctx, policy.UserLoginMustBeDomain))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return orgWriteModelToOrgIAMPolicy(addedPolicy), nil
}

func (r *CommandSide) addOrgIAMPolicy(ctx context.Context, orgAgg *org.Aggregate, addedPolicy *ORGOrgIAMPolicyWriteModel, policy *domain.OrgIAMPolicy) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "ORG-5M0ds", "Errors.Org.OrgIAMPolicy.AlreadyExists")
	}
	orgAgg.PushEvents(org.NewOrgIAMPolicyAddedEvent(ctx, policy.UserLoginMustBeDomain))
	return nil
}

func (r *CommandSide) ChangeOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	existingPolicy, err := r.orgIAMPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-2N9sd", "Errors.Org.OrgIAMPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-3M9ds", "Errors.Org.LabelPolicy.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return orgWriteModelToOrgIAMPolicy(existingPolicy), nil
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
