package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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

func (e *PasswordAgePolicyAddedEvent) Data() interface{} {
	return e
}

func (e *PasswordAgePolicyAddedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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

func PasswordAgePolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-T3mGp", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays *uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     *uint64 `json:"maxAgeDays,omitempty"`
}

func (e *PasswordAgePolicyChangedEvent) Data() interface{} {
	return e
}

func (e *PasswordAgePolicyChangedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewPasswordAgePolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PasswordAgePolicyChanges,
) (*PasswordAgePolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-DAgt5", "Errors.NoChangesFound")
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

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-PqaVq", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordAgePolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *PasswordAgePolicyRemovedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewPasswordAgePolicyRemovedEvent(base *eventstore.BaseEvent) *PasswordAgePolicyRemovedEvent {
	return &PasswordAgePolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordAgePolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PasswordAgePolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
