package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultPasswordAgePolicy(ctx context.Context, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	addedPolicy := NewIAMPasswordAgePolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	err := r.addDefaultPasswordAgePolicy(ctx, nil, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(&addedPolicy.PasswordAgePolicyWriteModel), nil
}

func (r *CommandSide) addDefaultPasswordAgePolicy(ctx context.Context, iamAgg *iam_repo.Aggregate, addedPolicy *IAMPasswordAgePolicyWriteModel, policy *domain.PasswordAgePolicy) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.PasswordAgePolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewPasswordAgePolicyAddedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays))

	return nil
}

func (r *CommandSide) ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	existingPolicy, err := r.defaultPasswordAgePolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0oPew", "Errors.IAM.PasswordAgePolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordAgePolicyWriteModel.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordAgePolicy(&existingPolicy.PasswordAgePolicyWriteModel), nil
}

func (r *CommandSide) defaultPasswordAgePolicyWriteModelByID(ctx context.Context) (policy *IAMPasswordAgePolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordAgePolicyWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
