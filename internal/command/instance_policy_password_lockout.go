package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultLockoutPolicy(ctx context.Context, maxPasswordAttempts, maxOTPAttempts uint64, showLockoutFailure bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultLockoutPolicy(
		instanceAgg,
		maxPasswordAttempts,
		maxOTPAttempts,
		showLockoutFailure,
	))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeDefaultLockoutPolicy(ctx context.Context, policy *domain.LockoutPolicy) (*domain.LockoutPolicy, error) {
	existingPolicy, err := c.defaultLockoutPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-0oPew", "Errors.IAM.LockoutPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LockoutPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		instanceAgg,
		policy.MaxPasswordAttempts,
		policy.MaxOTPAttempts,
		policy.ShowLockOutFailures,
	)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-0psjF", "Errors.IAM.LockoutPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLockoutPolicy(&existingPolicy.LockoutPolicyWriteModel), nil
}

func (c *Commands) defaultLockoutPolicyWriteModelByID(ctx context.Context) (policy *InstanceLockoutPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceLockoutPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func prepareAddDefaultLockoutPolicy(
	a *instance.Aggregate,
	maxPasswordAttempts,
	maxOTPAttempts uint64,
	showLockoutFailure bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceLockoutPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-0olDf", "Errors.Instance.LockoutPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewLockoutPolicyAddedEvent(ctx, &a.Aggregate, maxPasswordAttempts, maxOTPAttempts, showLockoutFailure),
			}, nil
		}, nil
	}
}
