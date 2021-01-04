package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetDefaultPasswordComplexityPolicy(ctx context.Context, aggregateID string) (*iam_model.PasswordComplexityPolicy, error) {
	policyWriteModel := NewIAMPasswordComplexityPolicyWriteModel(aggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToPasswordComplexityPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	addedPolicy := NewIAMPasswordComplexityPolicyWriteModel(policy.AggregateID)
	iamAgg, err := r.addDefaultPasswordComplexityPolicy(ctx, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordComplexityPolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultPasswordComplexityPolicy(ctx context.Context, addedPolicy *IAMPasswordComplexityPolicyWriteModel, policy *iam_model.PasswordComplexityPolicy) (*iam_repo.Aggregate, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.PasswordComplexityPolicy.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.PasswordComplexityPolicyWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewPasswordComplexityPolicyAddedEvent(ctx, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol))

	return iamAgg, nil
}

func (r *CommandSide) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	existingPolicy, err := r.defaultPasswordComplexityPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordComplexityPolicyWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordComplexityPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultPasswordComplexityPolicyWriteModelByID(ctx context.Context, iamID string) (policy *IAMPasswordComplexityPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordComplexityPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
