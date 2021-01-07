package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*domain.PasswordComplexityPolicy, error) {
	policyWriteModel := NewIAMPasswordComplexityPolicyWriteModel(r.iamID)
	err := r.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	policy := writeModelToPasswordComplexityPolicy(policyWriteModel)
	policy.Default = true
	return policy, nil
}

func (r *CommandSide) AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *domain.PasswordComplexityPolicy) (*domain.PasswordComplexityPolicy, error) {
	policy.AggregateID = r.iamID
	addedPolicy := NewIAMPasswordComplexityPolicyWriteModel(policy.AggregateID)
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	err := r.addDefaultPasswordComplexityPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordComplexityPolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultPasswordComplexityPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMPasswordComplexityPolicyWriteModel, policy *domain.PasswordComplexityPolicy) error {
	if err := policy.IsValid(); err != nil {
		return err
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.PasswordComplexityPolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewPasswordComplexityPolicyAddedEvent(ctx, policy.MinLength, policy.HasLowercase, policy.HasUppercase, policy.HasNumber, policy.HasSymbol))

	return nil
}

func (r *CommandSide) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *domain.PasswordComplexityPolicy) (*domain.PasswordComplexityPolicy, error) {
	policy.AggregateID = r.iamID
	if err := policy.IsValid(); err != nil {
		return nil, err
	}

	existingPolicy, err := r.defaultPasswordComplexityPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
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
