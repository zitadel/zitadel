package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	PasswordLockoutPolicyAddedEventType   = "policy.password.lockout.added"
	PasswordLockoutPolicyChangedEventType = "policy.password.lockout.changed"
	PasswordLockoutPolicyRemovedEventType = "policy.password.lockout.removed"
)

type PasswordLockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         uint64 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

func (e *PasswordLockoutPolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyAddedEvent(
	base *eventstore.BaseEvent,
	maxAttempts uint64,
	showLockOutFailures bool,
) *PasswordLockoutPolicyAddedEvent {

	return &PasswordLockoutPolicyAddedEvent{
		BaseEvent:           *base,
		MaxAttempts:         maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

func PasswordLockoutPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordLockoutPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-8XiVd", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordLockoutPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         *uint64 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures *bool   `json:"showLockOutFailures,omitempty"`
}

func (e *PasswordLockoutPolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PasswordLockoutPolicyChanges,
) *PasswordLockoutPolicyChangedEvent {
	changeEvent := &PasswordLockoutPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent
}

type PasswordLockoutPolicyChanges func(*PasswordLockoutPolicyChangedEvent)

func ChangeMaxAttempts(maxAttempts uint64) func(*PasswordLockoutPolicyChangedEvent) {
	return func(e *PasswordLockoutPolicyChangedEvent) {
		e.MaxAttempts = &maxAttempts
	}
}

func ChangeShowLockOutFailures(showLockOutFailures bool) func(*PasswordLockoutPolicyChangedEvent) {
	return func(e *PasswordLockoutPolicyChangedEvent) {
		e.ShowLockOutFailures = &showLockOutFailures
	}
}

func PasswordLockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordLockoutPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-lWGRc", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordLockoutPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordLockoutPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordLockoutPolicyRemovedEvent(base *eventstore.BaseEvent) *PasswordLockoutPolicyRemovedEvent {
	return &PasswordLockoutPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordLockoutPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PasswordLockoutPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
