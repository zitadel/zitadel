package policy

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	LockoutPolicyAddedEventType   = "policy.lockout.added"
	LockoutPolicyChangedEventType = "policy.lockout.changed"
	LockoutPolicyRemovedEventType = "policy.lockout.removed"
)

type LockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      uint64 `json:"maxOTPAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

func (e *LockoutPolicyAddedEvent) Payload() interface{} {
	return e
}

func (e *LockoutPolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewLockoutPolicyAddedEvent(
	base *eventstore.BaseEvent,
	maxPasswordAttempts,
	maxOTPAttempts uint64,
	showLockOutFailures bool,
) *LockoutPolicyAddedEvent {

	return &LockoutPolicyAddedEvent{
		BaseEvent:           *base,
		MaxPasswordAttempts: maxPasswordAttempts,
		MaxOTPAttempts:      maxOTPAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

func LockoutPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &LockoutPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-8XiVd", "unable to unmarshal policy")
	}

	return e, nil
}

type LockoutPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxPasswordAttempts *uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      *uint64 `json:"maxOTPAttempts,omitempty"`
	ShowLockOutFailures *bool   `json:"showLockOutFailures,omitempty"`
}

func (e *LockoutPolicyChangedEvent) Payload() interface{} {
	return e
}

func (e *LockoutPolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewLockoutPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LockoutPolicyChanges,
) (*LockoutPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-sdgh6", "Errors.NoChangesFound")
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

func ChangeMaxPasswordAttempts(maxAttempts uint64) func(*LockoutPolicyChangedEvent) {
	return func(e *LockoutPolicyChangedEvent) {
		e.MaxPasswordAttempts = &maxAttempts
	}
}

func ChangeMaxOTPAttempts(maxAttempts uint64) func(*LockoutPolicyChangedEvent) {
	return func(e *LockoutPolicyChangedEvent) {
		e.MaxOTPAttempts = &maxAttempts
	}
}

func ChangeShowLockOutFailures(showLockOutFailures bool) func(*LockoutPolicyChangedEvent) {
	return func(e *LockoutPolicyChangedEvent) {
		e.ShowLockOutFailures = &showLockOutFailures
	}
}

func LockoutPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &LockoutPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-lWGRc", "unable to unmarshal policy")
	}

	return e, nil
}

type LockoutPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LockoutPolicyRemovedEvent) Payload() interface{} {
	return nil
}

func (e *LockoutPolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewLockoutPolicyRemovedEvent(base *eventstore.BaseEvent) *LockoutPolicyRemovedEvent {
	return &LockoutPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LockoutPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &LockoutPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
