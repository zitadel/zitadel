package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultNotificationPolicy(ctx context.Context, passwordChange bool) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
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

func (c *Commands) ChangeDefaultNotificationPolicy(ctx context.Context, policy *domain.NotificationPolicy) (*domain.NotificationPolicy, error) {
	existingPolicy, err := c.defaultNotificationPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-x891na", "Errors.IAM.NotificationPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.NotificationPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, instanceAgg, policy.PasswordChange)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-29x02n", "Errors.IAM.NotificationPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToNotificationPolicy(&existingPolicy.NotificationPolicyWriteModel), nil
}

func (c *Commands) defaultNotificationPolicyWriteModelByID(ctx context.Context) (policy *InstanceNotificationPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceNotificationPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
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
				return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-xpo1bj", "Errors.Instance.NotificationPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewNotificationPolicyAddedEvent(ctx, &a.Aggregate, passwordChange),
			}, nil
		}, nil
	}
}

func ExistsDefaultNotificationPolicy(ctx context.Context, filter preparation.FilterToQueryReducer, instanceID string) (bool, error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderAsc().
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(instanceID).
		EventTypes(
			instance.NotificationPolicyAddedEventType,
		).Builder())
	if err != nil {
		return false, err
	}

	if len(events) > 0 {
		return true, nil
	}
	return false, nil
}
