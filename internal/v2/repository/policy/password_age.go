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

type PasswordAgePolicyReadModel struct {
	eventstore.ReadModel

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

	ExpireWarnDays int `json:"expireWarnDays,omitempty"`
	MaxAgeDays     int `json:"maxAgeDays,omitempty"`
}

func (e *PasswordAgePolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordAgePolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordAgePolicyAggregate,
) *PasswordAgePolicyChangedEvent {

	e := &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
	}

	if current.ExpireWarnDays != changed.ExpireWarnDays {
		e.ExpireWarnDays = changed.ExpireWarnDays
	}
	if current.MaxAgeDays != changed.MaxAgeDays {
		e.MaxAgeDays = changed.ExpireWarnDays
	}

	return e
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
			PasswordAgePolicyRemovedEventType,
		),
	}
}
