package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicy, error) {
	policy := NewORGOrgIAMPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	if policy.IsActive {
		return orgWriteModelToOrgIAMPolicy(policy), nil
	}
	return r.GetDefaultOrgIAMPolicy(ctx, r.iamID)
}

func (r *CommandSide) AddOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	addedPolicy := NewORGOrgIAMPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
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

func (r *CommandSide) ChangeOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	existingPolicy, err := r.orgIAMPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.IsActive {
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

func (r *CommandSide) orgIAMPolicyWriteModelByID(ctx context.Context, iamID string) (policy *ORGOrgIAMPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewORGOrgIAMPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
