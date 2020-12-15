package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
)

func (r *CommandSide) AddDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	addedPolicy := password_age.NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-6L0pd", "Errors.IAM.PasswordAgePolicy.AlreadyExists")
	}

	iamAgg := AggregateFromWriteModel(&addedPolicy.WriteModel.WriteModel).
		PushPasswordAgePolicyAddedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays)

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	existingPolicy, err := r.defaultPasswordAgePolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := AggregateFromWriteModel(&existingPolicy.WriteModel.WriteModel).
		PushPasswordAgePolicyChangedFromExisting(ctx, existingPolicy, policy.ExpireWarnDays, policy.MaxAgeDays)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(existingPolicy), nil
}

func (r *CommandSide) defaultPasswordAgePolicyWriteModelByID(ctx context.Context, iamID string) (policy *password_age.WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := password_age.NewWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
