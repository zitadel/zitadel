package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordAgePolicyAddedEventType   = "policy.password.age.added"
	PasswordAgePolicyChangedEventType = "policy.password.age.changed"
	PasswordAgePolicyRemovedEventType = "policy.password.age.removed"
)

type PasswordAgePolicyAggregate struct {
	eventstore.Aggregate

	ExpireWarnDays int
	MaxAgeDays     int
}

type PasswordAgePolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays int `json:"expireWarnDays"`
	MaxAgeDays     int `json:"maxAgeDays"`
}

func (e *PasswordAgePolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays int,
) *PasswordAgePolicyAddedEvent {

	return &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyAddedEventType,
		),
		ExpireWarnDays: expireWarnDays,
		MaxAgeDays:     maxAgeDays,
	}
}

type PasswordAgePolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	current *PasswordAgePolicyAggregate
	changed *PasswordAgePolicyAggregate
}

func (e *PasswordAgePolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyChangedEvent) Data() interface{} {
	changes := map[string]interface{}{}

	if e.current.ExpireWarnDays != e.changed.ExpireWarnDays {
		changes["expireWarnDays"] = e.changed.ExpireWarnDays
	}
	if e.current.MaxAgeDays != e.changed.MaxAgeDays {
		changes["maxAgeDays"] = e.changed.ExpireWarnDays
	}

	return changes
}

func NewPasswordAgePolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordAgePolicyAggregate,
) *PasswordAgePolicyChangedEvent {

	return &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
		current: current,
		changed: changed,
	}
}

type PasswordAgePolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordAgePolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordAgePolicyRemovedEvent(
	ctx context.Context,
	current,
	changed *PasswordAgePolicyRemovedEvent,
) *PasswordAgePolicyChangedEvent {

	return &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
	}
}
