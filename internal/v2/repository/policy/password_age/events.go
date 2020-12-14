package password_age

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

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays uint64 `json:"expireWarnDays"`
	MaxAgeDays     uint64 `json:"maxAgeDays"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	expireWarnDays,
	maxAgeDays uint64,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent:      *base,
		ExpireWarnDays: expireWarnDays,
		MaxAgeDays:     maxAgeDays,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-T3mGp", "unable to unmarshal policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	expireWarnDays,
	maxAgeDays uint64,
) *ChangedEvent {

	e := &ChangedEvent{
		BaseEvent: *base,
	}

	if current.ExpireWarnDays != expireWarnDays {
		e.ExpireWarnDays = expireWarnDays
	}
	if current.MaxAgeDays != maxAgeDays {
		e.MaxAgeDays = maxAgeDays
	}

	return e
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-PqaVq", "unable to unmarshal policy")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *base,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-02878", "unable to unmarshal policy")
	}

	return e, nil
}
