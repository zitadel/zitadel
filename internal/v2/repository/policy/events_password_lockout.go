package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordLockoutPolicyAddedEventType = "policy.password.lockout.added"
)

type PasswordLockoutAggregate struct {
	eventstore.Aggregate

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
			LabelPolicyAddedEventType,
		),
		MaxAttempts:         maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

type PasswordLockoutPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	current *PasswordLockoutAggregate
	changed *PasswordLockoutAggregate
}

func (e *PasswordLockoutPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyChangedEvent) Data() interface{} {
	changes := map[string]interface{}{}

	if e.current.MaxAttempts != e.changed.MaxAttempts {
		changes["maxAttempts"] = e.changed.MaxAttempts
	}
	if e.current.ShowLockOutFailures != e.changed.ShowLockOutFailures {
		changes["showLockOutFailures"] = e.changed.ShowLockOutFailures
	}

	return changes
}

func NewPasswordLockoutPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordLockoutAggregate,
) *PasswordLockoutPolicyChangedEvent {

	return &PasswordLockoutPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyAddedEventType,
		),
		current: current,
		changed: changed,
	}
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
			LabelPolicyAddedEventType,
		),
	}
}
