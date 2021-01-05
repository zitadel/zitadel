package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultPasswordAgePolicy(ctx context.Context, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	policy.AggregateID = r.iamID
	addedPolicy := NewIAMPasswordAgePolicyWriteModel(policy.AggregateID)
	iamAgg, err := r.addDefaultPasswordAgePolicy(ctx, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(addedPolicy), nil
}

func (r *CommandSide) addDefaultPasswordAgePolicy(ctx context.Context, addedPolicy *IAMPasswordAgePolicyWriteModel, policy *iam_model.PasswordAgePolicy) (*iam_repo.Aggregate, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.PasswordAgePolicy.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.PasswordAgePolicyWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewPasswordAgePolicyAddedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays))

	return iamAgg, nil
}

func (r *CommandSide) ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	policy.AggregateID = r.iamID
	existingPolicy, err := r.defaultPasswordAgePolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}
	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays)
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
