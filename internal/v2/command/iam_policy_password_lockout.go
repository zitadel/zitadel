package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	policy.AggregateID = r.iamID
	addedPolicy := NewIAMPasswordLockoutPolicyWriteModel(policy.AggregateID)
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	err := r.addDefaultPasswordLockoutPolicy(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultPasswordLockoutPolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMPasswordLockoutPolicyWriteModel, policy *domain.PasswordLockoutPolicy) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.IsActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-0olDf", "Errors.IAM.PasswordLockoutPolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewPasswordLockoutPolicyAddedEvent(ctx, policy.MaxAttempts, policy.ShowLockOutFailures))

	return nil
}

func (r *CommandSide) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	policy.AggregateID = r.iamID
	existingPolicy, err := r.defaultPasswordLockoutPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0oPew", "Errors.IAM.PasswordLockoutPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.MaxAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.PasswordLockoutPolicy.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordLockoutPolicyWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultPasswordLockoutPolicyWriteModelByID(ctx context.Context, iamID string) (policy *IAMPasswordLockoutPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordLockoutPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
