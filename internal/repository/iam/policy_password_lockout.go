package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	PasswordLockoutPolicyAddedEventType   = iamEventTypePrefix + policy.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = iamEventTypePrefix + policy.PasswordLockoutPolicyChangedEventType
)

type PasswordLockoutPolicyAddedEvent struct {
	policy.PasswordLockoutPolicyAddedEvent
}

func NewPasswordLockoutPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool,
) *PasswordLockoutPolicyAddedEvent {
	return &PasswordLockoutPolicyAddedEvent{
		PasswordLockoutPolicyAddedEvent: *policy.NewPasswordLockoutPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PasswordLockoutPolicyAddedEventType),
			maxAttempts,
			showLockoutFailure),
	}
}

func PasswordLockoutPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyAddedEvent{PasswordLockoutPolicyAddedEvent: *e.(*policy.PasswordLockoutPolicyAddedEvent)}, nil
}

type PasswordLockoutPolicyChangedEvent struct {
	policy.PasswordLockoutPolicyChangedEvent
}

func NewPasswordLockoutPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.PasswordLockoutPolicyChanges,
) (*PasswordLockoutPolicyChangedEvent, error) {
	changedEvent, err := policy.NewPasswordLockoutPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordLockoutPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PasswordLockoutPolicyChangedEvent{PasswordLockoutPolicyChangedEvent: *changedEvent}, nil
}

func PasswordLockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyChangedEvent{PasswordLockoutPolicyChangedEvent: *e.(*policy.PasswordLockoutPolicyChangedEvent)}, nil
}
