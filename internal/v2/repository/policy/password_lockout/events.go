package password_lockout

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

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         uint64 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	maxAttempts uint64,
	showLockOutFailures bool,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent:           *base,
		MaxAttempts:         maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-8XiVd", "unable to unmarshal policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         uint64 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	maxAttempts uint64,
	showLockOutFailures bool,
) *ChangedEvent {

	e := &ChangedEvent{
		BaseEvent: *base,
	}

	if current.MaxAttempts != maxAttempts {
		e.MaxAttempts = maxAttempts
	}
	if current.ShowLockOutFailures != showLockOutFailures {
		e.ShowLockOutFailures = showLockOutFailures
	}

	return e
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-lWGRc", "unable to unmarshal policy")
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
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
