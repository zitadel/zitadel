package password_complexity

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	PasswordComplexityPolicyAddedEventType   = "policy.password.complexity.added"
	PasswordComplexityPolicyChangedEventType = "policy.password.complexity.changed"
	PasswordComplexityPolicyRemovedEventType = "policy.password.complexity.removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint64 `json:"minLength,omitempty"`
	HasLowercase bool   `json:"hasLowercase"`
	HasUpperCase bool   `json:"hasUppercase"`
	HasNumber    bool   `json:"hasNumber"`
	HasSymbol    bool   `json:"hasSymbol"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	minLength uint64,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent:    *base,
		MinLength:    minLength,
		HasLowercase: hasLowerCase,
		HasNumber:    hasNumber,
		HasSymbol:    hasSymbol,
		HasUpperCase: hasUpperCase,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-wYxlM", "unable to unmarshal policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint64 `json:"minLength"`
	HasLowercase bool   `json:"hasLowercase"`
	HasUpperCase bool   `json:"hasUppercase"`
	HasNumber    bool   `json:"hasNumber"`
	HasSymbol    bool   `json:"hasSymbol"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	minLength uint64,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
) *ChangedEvent {

	e := &ChangedEvent{
		BaseEvent: *base,
	}

	if current.MinLength != minLength {
		e.MinLength = minLength
	}
	if current.HasLowercase != hasLowerCase {
		e.HasLowercase = hasLowerCase
	}
	if current.HasUpperCase != hasUpperCase {
		e.HasUpperCase = hasUpperCase
	}
	if current.HasNumber != hasNumber {
		e.HasNumber = hasNumber
	}
	if current.HasSymbol != hasSymbol {
		e.HasSymbol = hasSymbol
	}

	return e
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-zBGB0", "unable to unmarshal policy")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(base *eventstore.BaseEvent) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *base,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
