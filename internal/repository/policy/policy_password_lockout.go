package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func (e *PasswordLockoutPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *PasswordLockoutPolicyAddedEvent) Assets() []*eventstore.Asset {
	return nil
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

func (e *PasswordLockoutPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *PasswordLockoutPolicyChangedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewPasswordLockoutPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PasswordLockoutPolicyChanges,
) (*PasswordLockoutPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-sdgh6", "Errors.NoChangesFound")
	}
	changeEvent := &PasswordLockoutPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
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

func (e *PasswordLockoutPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *PasswordLockoutPolicyRemovedEvent) Assets() []*eventstore.Asset {
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
