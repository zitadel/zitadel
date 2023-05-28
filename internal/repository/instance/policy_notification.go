package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	NotificationPolicyAddedEventType   = instanceEventTypePrefix + policy.NotificationPolicyAddedEventType
	NotificationPolicyChangedEventType = instanceEventTypePrefix + policy.NotificationPolicyChangedEventType
)

type NotificationPolicyAddedEvent struct {
	policy.NotificationPolicyAddedEvent
}

func NewNotificationPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	passwordChange bool,
) *NotificationPolicyAddedEvent {
	return &NotificationPolicyAddedEvent{
		NotificationPolicyAddedEvent: *policy.NewNotificationPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				NotificationPolicyAddedEventType),
			passwordChange),
	}
}

func NotificationPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.NotificationPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &NotificationPolicyAddedEvent{NotificationPolicyAddedEvent: *e.(*policy.NotificationPolicyAddedEvent)}, nil
}

type NotificationPolicyChangedEvent struct {
	policy.NotificationPolicyChangedEvent
}

func NewNotificationPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.NotificationPolicyChanges,
) (*NotificationPolicyChangedEvent, error) {
	changedEvent, err := policy.NewNotificationPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			NotificationPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &NotificationPolicyChangedEvent{NotificationPolicyChangedEvent: *changedEvent}, nil
}

func NotificationPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.NotificationPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &NotificationPolicyChangedEvent{NotificationPolicyChangedEvent: *e.(*policy.NotificationPolicyChangedEvent)}, nil
}
