package command

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/command/preparation"
	"github.com/zitadel/zitadel/v2/internal/domain"
	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/org"
)

func (c *Commands) AddNotificationPolicy(ctx context.Context, resourceOwner string, passwordChange bool) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x801sk2i", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddNotificationPolicy(orgAgg, passwordChange))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareAddNotificationPolicy(
	a *org.Aggregate,
	passwordChange bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgNotificationPolicyWriteModel(a.Aggregate.ID)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, caos_errs.ThrowAlreadyExists(nil, "Org-xa08n2", "Errors.Org.NotificationPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				org.NewNotificationPolicyAddedEvent(ctx, &a.Aggregate, passwordChange),
			}, nil
		}, nil
	}
}

func (c *Commands) ChangeNotificationPolicy(ctx context.Context, resourceOwner string, passwordChange bool) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x091n1g", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeNotificationPolicy(orgAgg, passwordChange))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareChangeNotificationPolicy(
	a *org.Aggregate,
	passwordChange bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgNotificationPolicyWriteModel(a.Aggregate.ID)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}

			if writeModel.State == domain.PolicyStateUnspecified || writeModel.State == domain.PolicyStateRemoved {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-x029n3", "Errors.Org.NotificationPolicy.NotFound")
			}
			change, hasChanged := writeModel.NewChangedEvent(ctx, &a.Aggregate, passwordChange)
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-ioqnxz", "Errors.Org.NotificationPolicy.NotChanged")
			}
			return []eventstore.Command{
				change,
			}, nil
		}, nil
	}
}

func (c *Commands) RemoveNotificationPolicy(ctx context.Context, resourceOwner string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x89ns2", "Errors.ResourceOwnerMissing")
	}
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareRemoveNotificationPolicy(orgAgg))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareRemoveNotificationPolicy(
	a *org.Aggregate,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgNotificationPolicyWriteModel(a.Aggregate.ID)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}

			if writeModel.State == domain.PolicyStateUnspecified || writeModel.State == domain.PolicyStateRemoved {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-x029n1s", "Errors.Org.NotificationPolicy.NotFound")
			}
			return []eventstore.Command{
				org.NewNotificationPolicyRemovedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}
