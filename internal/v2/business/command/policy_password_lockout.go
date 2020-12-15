package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
)

func (r *CommandSide) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	addedPolicy := password_lockout.NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0olDf", "Errors.IAM.PasswordLockoutPolicy.AlreadyExists")
	}

	iamAgg := AggregateFromWriteModel(&addedPolicy.WriteModel.WriteModel).
		PushPasswordLockoutPolicyAddedEvent(ctx, policy.MaxAttempts, policy.ShowLockOutFailures)

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	existingPolicy, err := r.defaultPasswordLockoutPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := AggregateFromWriteModel(&existingPolicy.WriteModel.WriteModel).
		PushPasswordLockoutPolicyChangedFromExisting(ctx, existingPolicy, policy.MaxAttempts, policy.ShowLockOutFailures)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultPasswordLockoutPolicyWriteModelByID(ctx context.Context, iamID string) (policy *password_lockout.WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := password_lockout.NewWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
