package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	addedPolicy := NewIAMPasswordAgePolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.PasswordAgePolicy.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.PasswordAgePolicyWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewPasswordAgePolicyAddedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays))

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
	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(policy.ExpireWarnDays, policy.MaxAgeDays)
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordAgePolicyWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(existingPolicy), nil
}

func (r *CommandSide) defaultPasswordAgePolicyWriteModelByID(ctx context.Context, iamID string) (policy *IAMPasswordAgePolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordAgePolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
