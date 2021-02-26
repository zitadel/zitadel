package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	addedPolicy := NewIAMPasswordLockoutPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.WriteModel)
	event, err := c.addDefaultPasswordLockoutPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(&addedPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (c *Commands) addDefaultPasswordLockoutPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMPasswordLockoutPolicyWriteModel, policy *domain.PasswordLockoutPolicy) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0olDf", "Errors.IAM.PasswordLockoutPolicy.AlreadyExists")
	}

	return iam_repo.NewPasswordLockoutPolicyAddedEvent(ctx, iamAgg, policy.MaxAttempts, policy.ShowLockOutFailures), nil
}

func (c *Commands) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	existingPolicy, err := c.defaultPasswordLockoutPolicyWriteModelByID(ctx)
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

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordLockoutPolicy(&existingPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (c *Commands) defaultPasswordLockoutPolicyWriteModelByID(ctx context.Context) (policy *IAMPasswordLockoutPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMPasswordLockoutPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
