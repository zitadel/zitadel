package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	addedPolicy := NewIAMPasswordLockoutPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := r.addDefaultPasswordLockoutPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(&addedPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) addDefaultPasswordLockoutPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMPasswordLockoutPolicyWriteModel, policy *domain.PasswordLockoutPolicy) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0olDf", "Errors.IAM.PasswordLockoutPolicy.AlreadyExists")
	}

	return iam_repo.NewPasswordLockoutPolicyAddedEvent(ctx, iamAgg, policy.MaxAttempts, policy.ShowLockOutFailures), nil
}

func (r *CommandSide) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	existingPolicy, err := r.defaultPasswordLockoutPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0oPew", "Errors.IAM.PasswordLockoutPolicy.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.PasswordLockoutPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, iamAgg, policy.MaxAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.PasswordLockoutPolicy.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordLockoutPolicy(&existingPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) defaultPasswordLockoutPolicyWriteModelByID(ctx context.Context) (policy *IAMPasswordLockoutPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordLockoutPolicyWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
