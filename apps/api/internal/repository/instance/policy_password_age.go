package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	PasswordAgePolicyAddedEventType   = instanceEventTypePrefix + policy.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = instanceEventTypePrefix + policy.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyAddedEvent struct {
	policy.PasswordAgePolicyAddedEvent
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	expireWarnDays,
	maxAgeDays uint64,
) *PasswordAgePolicyAddedEvent {
	return &PasswordAgePolicyAddedEvent{
		PasswordAgePolicyAddedEvent: *policy.NewPasswordAgePolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PasswordAgePolicyAddedEventType),
			expireWarnDays,
			maxAgeDays),
	}
}

func PasswordAgePolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PasswordAgePolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyAddedEvent{PasswordAgePolicyAddedEvent: *e.(*policy.PasswordAgePolicyAddedEvent)}, nil
}

type PasswordAgePolicyChangedEvent struct {
	policy.PasswordAgePolicyChangedEvent
}

func NewPasswordAgePolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.PasswordAgePolicyChanges,
) (*PasswordAgePolicyChangedEvent, error) {
	changedEvent, err := policy.NewPasswordAgePolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordAgePolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *changedEvent}, nil
}

func PasswordAgePolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*policy.PasswordAgePolicyChangedEvent)}, nil
}
