package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

const (
	PasswordComplexityPolicyAddedEventType   = iamEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = iamEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
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
) *PasswordComplexityPolicyChangedEvent {
	return &PasswordComplexityPolicyChangedEvent{
		PasswordComplexityPolicyChangedEvent: *policy.NewPasswordComplexityPolicyChangedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyAddedEventType),
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
