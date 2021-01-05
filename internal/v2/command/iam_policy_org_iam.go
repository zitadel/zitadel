package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetDefaultOrgIAMPolicy(ctx context.Context) (*domain.OrgIAMPolicy, error) {
	policyWriteModel := NewIAMOrgIAMPolicyWriteModel(r.iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToOrgIAMPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	policy.AggregateID = r.iamID
	addedPolicy := NewIAMOrgIAMPolicyWriteModel(policy.AggregateID)
	iamAgg, err := r.addDefaultOrgIAMPolicy(ctx, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToOrgIAMPolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultOrgIAMPolicy(ctx context.Context, addedPolicy *IAMOrgIAMPolicyWriteModel, policy *iam_model.OrgIAMPolicy) (*iam_repo.Aggregate, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.OrgIAMPolicy.AlreadyExists")
	}
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.PolicyOrgIAMWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewOrgIAMPolicyAddedEvent(ctx, policy.UserLoginMustBeDomain))

	return iamAgg, nil
}

func (r *CommandSide) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *domain.OrgIAMPolicy) (*domain.OrgIAMPolicy, error) {
	policy.AggregateID = r.iamID
	existingPolicy, err := r.defaultOrgIAMPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0Pl0d", "Errors.IAM.OrgIAMPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.UserLoginMustBeDomain)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PolicyOrgIAMWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToOrgIAMPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultOrgIAMPolicyWriteModelByID(ctx context.Context, iamID string) (policy *IAMOrgIAMPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMOrgIAMPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
