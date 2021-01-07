package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetOrgIAMPolicy(ctx context.Context, orgID string) (*domain.OrgIAMPolicy, error) {
	policy := NewORGOrgIAMPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return orgWriteModelToOrgIAMPolicy(policy), nil
	}
	return r.GetDefaultOrgIAMPolicy(ctx)
}

func (r *CommandSide) AddOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	addedPolicy := NewORGOrgIAMPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-5M0ds", "Errors.Org.OrgIAMPolicy.AlreadyExists")
	}
	orgAgg := ORGAggregateFromWriteModel(&addedPolicy.PolicyOrgIAMWriteModel.WriteModel)
	orgAgg.PushEvents(iam_repo.NewOrgIAMPolicyAddedEvent(ctx, policy.UserLoginMustBeDomain))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return orgWriteModelToOrgIAMPolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	existingPolicy, err := r.orgIAMPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-2N9sd", "Errors.Org.OrgIAMPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-3M9ds", "Errors.Org.LabelPolicy.NotChanged")
	}

	orgAgg := ORGAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return orgWriteModelToOrgIAMPolicy(existingPolicy), nil
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
