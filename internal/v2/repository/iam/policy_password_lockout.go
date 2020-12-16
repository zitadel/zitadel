package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
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
	maxAttempts uint64,
	showLockoutFailure bool,
) *PasswordLockoutPolicyAddedEvent {
	return &PasswordLockoutPolicyAddedEvent{
		PasswordLockoutPolicyAddedEvent: *policy.NewPasswordLockoutPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordLockoutPolicyAddedEventType),
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

func PasswordLockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyChangedEvent{PasswordLockoutPolicyChangedEvent: *e.(*policy.PasswordLockoutPolicyChangedEvent)}, nil
}
