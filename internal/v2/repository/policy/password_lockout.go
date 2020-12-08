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

type PasswordLockoutPolicyAggregate struct {
	eventstore.Aggregate
}

type PasswordLockoutPolicyReadModel struct {
	eventstore.ReadModel

	MaxAttempts         uint8
	ShowLockOutFailures bool
}

func (rm *PasswordLockoutPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		case *PasswordLockoutPolicyChangedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		}
	}
	return rm.ReadModel.Reduce()
}

type PasswordLockoutPolicyWriteModel struct {
	eventstore.WriteModel

	MaxAttempts         uint8
	ShowLockOutFailures bool
}

func (wm *PasswordLockoutPolicyWriteModel) Reduce() error {
	return errors.ThrowUnimplemented(nil, "POLIC-3M0df", "reduce unimpelemnted")
}

type PasswordLockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         uint8 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool  `json:"showLockOutFailures"`
}

func (e *PasswordLockoutPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyAddedEvent(
	base *eventstore.BaseEvent,
	maxAttempts uint8,
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

	MaxAttempts         uint8 `json:"maxAttempts,omitempty"`
	ShowLockOutFailures bool  `json:"showLockOutFailures,omitempty"`
}

func (e *PasswordLockoutPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyChangedEvent(
	base *eventstore.BaseEvent,
	current *PasswordLockoutPolicyWriteModel,
	maxAttempts uint8,
	showLockOutFailures bool,
) *PasswordLockoutPolicyChangedEvent {

	e := &PasswordLockoutPolicyChangedEvent{
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

func (e *PasswordLockoutPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordLockoutPolicyRemovedEvent(
	base *eventstore.BaseEvent,
) *PasswordLockoutPolicyRemovedEvent {

	return &PasswordLockoutPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordLockoutPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PasswordLockoutPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
