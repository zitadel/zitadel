package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordComplexityPolicyAddedEventType   = orgEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = orgEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
	PasswordComplexityPolicyRemovedEventType = orgEventTypePrefix + policy.PasswordComplexityPolicyRemovedEventType
)

type PasswordComplexityPolicyAddedEvent struct {
	policy.PasswordComplexityPolicyAddedEvent
}

func NewPasswordComplexityPolicyAddedEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) *PasswordComplexityPolicyAddedEvent {
	return &PasswordComplexityPolicyAddedEvent{
		PasswordComplexityPolicyAddedEvent: *policy.NewPasswordComplexityPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyAddedEventType),
			minLength,
			hasLowercase,
			hasUppercase,
			hasNumber,
			hasSymbol),
	}
}

func PasswordComplexityPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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
	changes []policy.PasswordComplexityPolicyChanges,
) *PasswordComplexityPolicyChangedEvent {
	return &PasswordComplexityPolicyChangedEvent{
		PasswordComplexityPolicyChangedEvent: *policy.NewPasswordComplexityPolicyChangedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyChangedEventType),
			changes,
		),
	}
}

func PasswordComplexityPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordComplexityPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyChangedEvent{PasswordComplexityPolicyChangedEvent: *e.(*policy.PasswordComplexityPolicyChangedEvent)}, nil
}

type PasswordComplexityPolicyRemovedEvent struct {
	policy.PasswordComplexityPolicyRemovedEvent
}

func NewPasswordComplexityPolicyRemovedEvent(
	ctx context.Context,
) *PasswordComplexityPolicyRemovedEvent {
	return &PasswordComplexityPolicyRemovedEvent{
		PasswordComplexityPolicyRemovedEvent: *policy.NewPasswordComplexityPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyRemovedEventType),
		),
	}
}

func PasswordComplexityPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordComplexityPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyRemovedEvent{PasswordComplexityPolicyRemovedEvent: *e.(*policy.PasswordComplexityPolicyRemovedEvent)}, nil
}
