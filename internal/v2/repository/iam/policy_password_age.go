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
	policy.PassowordAgePolicyAddedEvent
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays uint64,
) *PasswordAgePolicyAddedEvent {
	return &PasswordAgePolicyAddedEvent{
		PassowordAgePolicyAddedEvent: *policy.NewPasswordAgePolicyAddedEvent(
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

	return &PasswordAgePolicyAddedEvent{PassowordAgePolicyAddedEvent: *e.(*policy.PassowordAgePolicyAddedEvent)}, nil
}

type PasswordAgePolicyChangedEvent struct {
	policy.PasswordAgePolicyChangedEvent
}

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*policy.PasswordAgePolicyChangedEvent)}, nil
}
