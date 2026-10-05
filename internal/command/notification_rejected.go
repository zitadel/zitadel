package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// NotificationRejected states on the aggregate that triggered a notification, that it was intentionally not sent.
// Notifications are triggered by events of users and sessions, which must still exist.
func (c *Commands) NotificationRejected(ctx context.Context, aggregate *eventstore.Aggregate, rejection channels.Rejection) error {
	if aggregate == nil || aggregate.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Rj3k1", "Errors.IDMissing")
	}
	var event eventstore.Command
	switch aggregate.Type {
	case user.AggregateType:
		// only the lifecycle events of the user are needed, not the whole write model
		resourceOwner, exists, err := existingUser(ctx, c.eventstore.Filter, aggregate.ID, aggregate.ResourceOwner, false) //nolint:staticcheck
		if err != nil {
			return err
		}
		if !exists {
			return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rj3k3", "Errors.User.NotFound")
		}
		event = user.NewNotificationRejectedEvent(ctx, &user.NewAggregate(aggregate.ID, resourceOwner).Aggregate, rejection)
	case session.AggregateType:
		sessionWriteModel := NewSessionWriteModel(aggregate.ID, aggregate.ResourceOwner)
		if err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
			return err
		}
		if err := sessionWriteModel.CheckIsActive(); err != nil {
			return err
		}
		event = session.NewNotificationRejectedEvent(ctx, &session.NewAggregate(aggregate.ID, sessionWriteModel.ResourceOwner).Aggregate, rejection)
	default:
		return zerrors.ThrowInvalidArgumentf(nil, "COMMAND-Rj3k2", "notifications of aggregate type %s cannot be rejected", aggregate.Type)
	}
	_, err := c.eventstore.Push(ctx, event)
	return err
}
