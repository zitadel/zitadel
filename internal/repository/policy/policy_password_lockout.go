package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	LockoutPolicyAddedEventType   = "policy.lockout.added"
	LockoutPolicyChangedEventType = "policy.lockout.changed"
	LockoutPolicyRemovedEventType = "policy.lockout.removed"
)

type LockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

func (e *LockoutPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *LockoutPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLockoutPolicyAddedEvent(
	base *eventstore.BaseEvent,
	maxAttempts uint64,
	showLockOutFailures bool,
) *LockoutPolicyAddedEvent {

	return &LockoutPolicyAddedEvent{
		BaseEvent:           *base,
		MaxPasswordAttempts: maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

func LockoutPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LockoutPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-8XiVd", "unable to unmarshal policy")
	}

	return e, nil
}

type LockoutPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxPasswordAttempts *uint64 `json:"maxPasswordAttempts,omitempty"`
	ShowLockOutFailures *bool   `json:"showLockOutFailures,omitempty"`
}

func (e *LockoutPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *LockoutPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLockoutPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LockoutPolicyChanges,
) (*LockoutPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-sdgh6", "Errors.NoChangesFound")
	}
	changeEvent := &LockoutPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type LockoutPolicyChanges func(*LockoutPolicyChangedEvent)

func ChangeMaxAttempts(maxAttempts uint64) func(*LockoutPolicyChangedEvent) {
	return func(e *LockoutPolicyChangedEvent) {
		e.MaxPasswordAttempts = &maxAttempts
	}
}

func ChangeShowLockOutFailures(showLockOutFailures bool) func(*LockoutPolicyChangedEvent) {
	return func(e *LockoutPolicyChangedEvent) {
		e.ShowLockOutFailures = &showLockOutFailures
	}
}

func LockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LockoutPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-lWGRc", "unable to unmarshal policy")
	}

	return e, nil
}

type LockoutPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LockoutPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *LockoutPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLockoutPolicyRemovedEvent(base *eventstore.BaseEvent) *LockoutPolicyRemovedEvent {
	return &LockoutPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LockoutPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LockoutPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
