package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	PasswordComplexityPolicyAddedEventType   = instanceEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = instanceEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
)

type PasswordComplexityPolicyAddedEvent struct {
	policy.PasswordComplexityPolicyAddedEvent
}

func NewPasswordComplexityPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) *PasswordComplexityPolicyAddedEvent {
	return &PasswordComplexityPolicyAddedEvent{
		PasswordComplexityPolicyAddedEvent: *policy.NewPasswordComplexityPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PasswordComplexityPolicyAddedEventType),
			minLength,
			hasLowercase,
			hasUppercase,
			hasNumber,
			hasSymbol),
	}
}

func PasswordComplexityPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PasswordComplexityPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyAddedEvent{PasswordComplexityPolicyAddedEvent: *e.(*policy.PasswordComplexityPolicyAddedEvent)}, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	policy.PasswordComplexityPolicyChangedEvent
}

func NewPasswordComplexityPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.PasswordComplexityPolicyChanges,
) (*PasswordComplexityPolicyChangedEvent, error) {
	changedEvent, err := policy.NewPasswordComplexityPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordComplexityPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PasswordComplexityPolicyChangedEvent{PasswordComplexityPolicyChangedEvent: *changedEvent}, nil
}

func PasswordComplexityPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PasswordComplexityPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyChangedEvent{PasswordComplexityPolicyChangedEvent: *e.(*policy.PasswordComplexityPolicyChangedEvent)}, nil
}
