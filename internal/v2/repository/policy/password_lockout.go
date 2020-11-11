package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordLockoutPolicyAddedEventType   = "policy.password.lockout.added"
	PasswordLockoutPolicyChangedEventType = "policy.password.lockout.changed"
	PasswordLockoutPolicyRemovedEventType = "policy.password.lockout.removed"
)

type PasswordLockoutPolicyAggregate struct {
	eventstore.Aggregate

	MaxAttempts         int
	ShowLockOutFailures bool
}

type PasswordLockoutPolicyReadModel struct {
	eventstore.ReadModel

	MaxAttempts         int
	ShowLockOutFailures bool
}

type PasswordLockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         int  `json:"maxAttempts"`
	ShowLockOutFailures bool `json:"showLockOutFailures"`
}

func (e *PasswordLockoutPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyAddedEvent(
	ctx context.Context,
	maxAttempts int,
	showLockOutFailures bool,
) *PasswordLockoutPolicyAddedEvent {

	return &PasswordLockoutPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordLockoutPolicyAddedEventType,
		),
		MaxAttempts:         maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

type PasswordLockoutPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         int  `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool `json:"showLockOutFailures,omitempty"`
}

func (e *PasswordLockoutPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordLockoutPolicyAggregate,
) *PasswordLockoutPolicyChangedEvent {

	e := &PasswordLockoutPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordLockoutPolicyChangedEventType,
		),
	}

	if current.MaxAttempts != changed.MaxAttempts {
		e.MaxAttempts = changed.MaxAttempts
	}
	if current.ShowLockOutFailures != changed.ShowLockOutFailures {
		e.ShowLockOutFailures = changed.ShowLockOutFailures
	}

	return e
}

type PasswordLockoutPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordLockoutPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordLockoutPolicyRemovedEvent(
	ctx context.Context,
) *PasswordLockoutPolicyRemovedEvent {

	return &PasswordLockoutPolicyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordLockoutPolicyRemovedEventType,
		),
	}
}
