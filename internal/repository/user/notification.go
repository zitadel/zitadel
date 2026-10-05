package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

const (
	NotificationRejectedType = userEventTypePrefix + "notification.rejected"
)

// NotificationRejectedEvent states that a notification to the user was intentionally not sent.
// It is not retried.
type NotificationRejectedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	channels.Rejection
}

func (e *NotificationRejectedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *NotificationRejectedEvent) Payload() any {
	return e
}

func (e *NotificationRejectedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotificationRejectedEvent(ctx context.Context, aggregate *eventstore.Aggregate, rejection channels.Rejection) *NotificationRejectedEvent {
	return &NotificationRejectedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			NotificationRejectedType,
		),
		Rejection: rejection,
	}
}
