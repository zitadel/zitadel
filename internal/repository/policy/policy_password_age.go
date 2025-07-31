package policy

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	PasswordAgePolicyAddedEventType   = "policy.password.age.added"
	PasswordAgePolicyChangedEventType = "policy.password.age.changed"
	PasswordAgePolicyRemovedEventType = "policy.password.age.removed"
)

type PasswordAgePolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

func (e *PasswordAgePolicyAddedEvent) Payload() any {
	return e
}

func (e *PasswordAgePolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordAgePolicyAddedEvent(
	base *eventstore.BaseEvent,
	expireWarnDays,
	maxAgeDays uint64,
) *PasswordAgePolicyAddedEvent {

	return &PasswordAgePolicyAddedEvent{
		BaseEvent:      *base,
		ExpireWarnDays: expireWarnDays,
		MaxAgeDays:     maxAgeDays,
	}
}

func PasswordAgePolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-T3mGp", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays *uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     *uint64 `json:"maxAgeDays,omitempty"`
}

func (e *PasswordAgePolicyChangedEvent) Payload() any {
	return e
}

func (e *PasswordAgePolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordAgePolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PasswordAgePolicyChanges,
) (*PasswordAgePolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-DAgt5", "Errors.NoChangesFound")
	}
	changeEvent := &PasswordAgePolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type PasswordAgePolicyChanges func(*PasswordAgePolicyChangedEvent)

func ChangeExpireWarnDays(expireWarnDay uint64) func(*PasswordAgePolicyChangedEvent) {
	return func(e *PasswordAgePolicyChangedEvent) {
		e.ExpireWarnDays = &expireWarnDay
	}
}

func ChangeMaxAgeDays(maxAgeDays uint64) func(*PasswordAgePolicyChangedEvent) {
	return func(e *PasswordAgePolicyChangedEvent) {
		e.MaxAgeDays = &maxAgeDays
	}
}

func PasswordAgePolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-PqaVq", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordAgePolicyRemovedEvent) Payload() any {
	return nil
}

func (e *PasswordAgePolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordAgePolicyRemovedEvent(base *eventstore.BaseEvent) *PasswordAgePolicyRemovedEvent {
	return &PasswordAgePolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordAgePolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PasswordAgePolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
