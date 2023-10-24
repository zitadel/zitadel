package policy

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	PasswordComplexityPolicyAddedEventType   = "policy.password.complexity.added"
	PasswordComplexityPolicyChangedEventType = "policy.password.complexity.changed"
	PasswordComplexityPolicyRemovedEventType = "policy.password.complexity.removed"
)

type PasswordComplexityPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint64 `json:"minLength,omitempty"`
	HasLowercase bool   `json:"hasLowercase,omitempty"`
	HasUppercase bool   `json:"hasUppercase,omitempty"`
	HasNumber    bool   `json:"hasNumber,omitempty"`
	HasSymbol    bool   `json:"hasSymbol,omitempty"`
}

func (e *PasswordComplexityPolicyAddedEvent) Payload() interface{} {
	return e
}

func (e *PasswordComplexityPolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordComplexityPolicyAddedEvent(
	base *eventstore.BaseEvent,
	minLength uint64,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
) *PasswordComplexityPolicyAddedEvent {
	return &PasswordComplexityPolicyAddedEvent{
		BaseEvent:    *base,
		MinLength:    minLength,
		HasLowercase: hasLowerCase,
		HasUppercase: hasUpperCase,
		HasNumber:    hasNumber,
		HasSymbol:    hasSymbol,
	}
}

func PasswordComplexityPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PasswordComplexityPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-wYxlM", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    *uint64 `json:"minLength,omitempty"`
	HasLowercase *bool   `json:"hasLowercase,omitempty"`
	HasUppercase *bool   `json:"hasUppercase,omitempty"`
	HasNumber    *bool   `json:"hasNumber,omitempty"`
	HasSymbol    *bool   `json:"hasSymbol,omitempty"`
}

func (e *PasswordComplexityPolicyChangedEvent) Payload() interface{} {
	return e
}

func (e *PasswordComplexityPolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordComplexityPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PasswordComplexityPolicyChanges,
) (*PasswordComplexityPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-Rdhu3", "Errors.NoChangesFound")
	}
	changeEvent := &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type PasswordComplexityPolicyChanges func(*PasswordComplexityPolicyChangedEvent)

func ChangeMinLength(minLength uint64) func(*PasswordComplexityPolicyChangedEvent) {
	return func(e *PasswordComplexityPolicyChangedEvent) {
		e.MinLength = &minLength
	}
}

func ChangeHasLowercase(hasLowercase bool) func(*PasswordComplexityPolicyChangedEvent) {
	return func(e *PasswordComplexityPolicyChangedEvent) {
		e.HasLowercase = &hasLowercase
	}
}

func ChangeHasUppercase(hasUppercase bool) func(*PasswordComplexityPolicyChangedEvent) {
	return func(e *PasswordComplexityPolicyChangedEvent) {
		e.HasUppercase = &hasUppercase
	}
}

func ChangeHasNumber(hasNumber bool) func(*PasswordComplexityPolicyChangedEvent) {
	return func(e *PasswordComplexityPolicyChangedEvent) {
		e.HasNumber = &hasNumber
	}
}

func ChangeHasSymbol(hasSymbol bool) func(*PasswordComplexityPolicyChangedEvent) {
	return func(e *PasswordComplexityPolicyChangedEvent) {
		e.HasSymbol = &hasSymbol
	}
}

func PasswordComplexityPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-zBGB0", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordComplexityPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordComplexityPolicyRemovedEvent) Payload() interface{} {
	return nil
}

func (e *PasswordComplexityPolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordComplexityPolicyRemovedEvent(base *eventstore.BaseEvent) *PasswordComplexityPolicyRemovedEvent {
	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordComplexityPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
