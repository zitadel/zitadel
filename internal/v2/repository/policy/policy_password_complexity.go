package policy

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

type PasswordComplexityPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint64 `json:"minLength,omitempty"`
	HasLowercase bool   `json:"hasLowercase,omitempty"`
	HasUpperCase bool   `json:"hasUppercase,omitempty"`
	HasNumber    bool   `json:"hasNumber,omitempty"`
	HasSymbol    bool   `json:"hasSymbol,omitempty"`
}

func (e *PasswordComplexityPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *PasswordComplexityPolicyAddedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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
		HasNumber:    hasNumber,
		HasSymbol:    hasSymbol,
		HasUpperCase: hasUpperCase,
	}
}

func PasswordComplexityPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordComplexityPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-wYxlM", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    *uint64 `json:"minLength,omitempty"`
	HasLowercase *bool   `json:"hasLowercase,omitempty"`
	HasUpperCase *bool   `json:"hasUppercase,omitempty"`
	HasNumber    *bool   `json:"hasNumber,omitempty"`
	HasSymbol    *bool   `json:"hasSymbol,omitempty"`
}

func (e *PasswordComplexityPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *PasswordComplexityPolicyChangedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewPasswordComplexityPolicyChangedEvent(
	base *eventstore.BaseEvent,
) *PasswordComplexityPolicyChangedEvent {
	return &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *base,
	}
}

func PasswordComplexityPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-zBGB0", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordComplexityPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordComplexityPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *PasswordComplexityPolicyRemovedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewPasswordComplexityPolicyRemovedEvent(base *eventstore.BaseEvent) *PasswordComplexityPolicyRemovedEvent {
	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PasswordComplexityPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
