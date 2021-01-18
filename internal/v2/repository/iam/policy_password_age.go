package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordAgePolicyAddedEventType   = iamEventTypePrefix + policy.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = iamEventTypePrefix + policy.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyAddedEvent struct {
	policy.PasswordAgePolicyAddedEvent
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays uint64,
) *PasswordAgePolicyAddedEvent {
	return &PasswordAgePolicyAddedEvent{
		PasswordAgePolicyAddedEvent: *policy.NewPasswordAgePolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordAgePolicyAddedEventType),
			expireWarnDays,
			maxAgeDays),
	}
}

func PasswordAgePolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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
	changes []policy.PasswordAgePolicyChanges,
) (*PasswordAgePolicyChangedEvent, error) {
	changedEvent, err := policy.NewPasswordAgePolicyChangedEvent(
		eventstore.NewBaseEventForPush(ctx, PasswordAgePolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *changedEvent}, nil
}

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*policy.PasswordAgePolicyChangedEvent)}, nil
}
