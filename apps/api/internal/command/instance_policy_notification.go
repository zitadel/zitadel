package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDefaultNotificationPolicy(ctx context.Context, resourceOwner string, passwordChange bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddDefaultNotificationPolicy(instanceAgg, passwordChange))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) ChangeDefaultNotificationPolicy(ctx context.Context, resourceOwner string, passwordChange bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeDefaultNotificationPolicy(instanceAgg, passwordChange))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareAddDefaultNotificationPolicy(
	a *instance.Aggregate,
	passwordChange bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceNotificationPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-xpo1bj", "Errors.Instance.NotificationPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewNotificationPolicyAddedEvent(ctx, &a.Aggregate, passwordChange),
			}, nil
		}, nil
	}
}

func prepareChangeDefaultNotificationPolicy(
	a *instance.Aggregate,
	passwordChange bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceNotificationPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}

			if writeModel.State == domain.PolicyStateUnspecified || writeModel.State == domain.PolicyStateRemoved {
				return nil, zerrors.ThrowNotFound(nil, "INSTANCE-x891na", "Errors.IAM.NotificationPolicy.NotFound")
			}
			change, hasChanged := writeModel.NewChangedEvent(ctx, &a.Aggregate, passwordChange)
			if !hasChanged {
				return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-29x02n", "Errors.IAM.NotificationPolicy.NotChanged")
			}
			return []eventstore.Command{
				change,
			}, nil
		}, nil
	}
}
