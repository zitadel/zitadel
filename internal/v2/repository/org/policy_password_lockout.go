package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordLockoutPolicyAddedEventType   = orgEventTypePrefix + policy.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = orgEventTypePrefix + policy.PasswordLockoutPolicyChangedEventType
	PasswordLockoutPolicyRemovedEventType = orgEventTypePrefix + policy.PasswordLockoutPolicyRemovedEventType
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

type PasswordLockoutPolicyRemovedEvent struct {
	policy.PasswordLockoutPolicyRemovedEvent
}

func NewPasswordLockoutPolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PasswordLockoutPolicyRemovedEvent {
	return &PasswordLockoutPolicyRemovedEvent{
		PasswordLockoutPolicyRemovedEvent: *policy.NewPasswordLockoutPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PasswordLockoutPolicyRemovedEventType),
		),
	}
}

func PasswordLockoutPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyRemovedEvent{PasswordLockoutPolicyRemovedEvent: *e.(*policy.PasswordLockoutPolicyRemovedEvent)}, nil
}
