package org

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordAgePolicyAddedEventType   = orgEventTypePrefix + policy.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = orgEventTypePrefix + policy.PasswordAgePolicyChangedEventType
	PasswordAgePolicyRemovedEventType = orgEventTypePrefix + policy.PasswordAgePolicyRemovedEventType
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

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*policy.PasswordAgePolicyChangedEvent)}, nil
}

type PasswordAgePolicyRemovedEvent struct {
	policy.PasswordAgePolicyRemovedEvent
}

func NewPasswordAgePolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PasswordAgePolicyRemovedEvent {
	return &PasswordAgePolicyRemovedEvent{
		PasswordAgePolicyRemovedEvent: *policy.NewPasswordAgePolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PasswordAgePolicyRemovedEventType),
		),
	}
}

func PasswordAgePolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyRemovedEvent{PasswordAgePolicyRemovedEvent: *e.(*policy.PasswordAgePolicyRemovedEvent)}, nil
}
