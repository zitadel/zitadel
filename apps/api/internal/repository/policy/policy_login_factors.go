package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	loginPolicySecondFactorPrefix           = loginPolicyPrefix + "secondfactor."
	LoginPolicySecondFactorAddedEventType   = loginPolicySecondFactorPrefix + "added"
	LoginPolicySecondFactorRemovedEventType = loginPolicySecondFactorPrefix + "removed"

	loginPolicyMultiFactorPrefix           = "policy.login.multifactor."
	LoginPolicyMultiFactorAddedEventType   = loginPolicyMultiFactorPrefix + "added"
	LoginPolicyMultiFactorRemovedEventType = loginPolicyMultiFactorPrefix + "removed"
)

type SecondFactorAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MFAType domain.SecondFactorType `json:"mfaType,omitempty"`
}

func NewSecondFactorAddedEvent(
	base *eventstore.BaseEvent,
	mfaType domain.SecondFactorType,
) *SecondFactorAddedEvent {
	return &SecondFactorAddedEvent{
		BaseEvent: *base,
		MFAType:   mfaType,
	}
}

func SecondFactorAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SecondFactorAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-Lp0dE", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *SecondFactorAddedEvent) Payload() interface{} {
	return e
}

func (e *SecondFactorAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SecondFactorRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	MFAType              domain.SecondFactorType `json:"mfaType"`
}

func NewSecondFactorRemovedEvent(
	base *eventstore.BaseEvent,
	mfaType domain.SecondFactorType,
) *SecondFactorRemovedEvent {
	return &SecondFactorRemovedEvent{
		BaseEvent: *base,
		MFAType:   mfaType,
	}
}

func SecondFactorRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SecondFactorRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-5M9gd", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *SecondFactorRemovedEvent) Payload() interface{} {
	return e
}

func (e *SecondFactorRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type MultiFactorAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MFAType domain.MultiFactorType `json:"mfaType"`
}

func NewMultiFactorAddedEvent(
	base *eventstore.BaseEvent,
	mfaType domain.MultiFactorType,
) *MultiFactorAddedEvent {
	return &MultiFactorAddedEvent{
		BaseEvent: *base,
		MFAType:   mfaType,
	}
}

func MultiFactorAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MultiFactorAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-5Ms90", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *MultiFactorAddedEvent) Payload() interface{} {
	return e
}

func (e *MultiFactorAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type MultiFactorRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	MFAType              domain.MultiFactorType `json:"mfaType"`
}

func NewMultiFactorRemovedEvent(
	base *eventstore.BaseEvent,
	mfaType domain.MultiFactorType,
) *MultiFactorRemovedEvent {
	return &MultiFactorRemovedEvent{
		BaseEvent: *base,
		MFAType:   mfaType,
	}
}

func MultiFactorRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MultiFactorRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-1N8sd", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *MultiFactorRemovedEvent) Payload() interface{} {
	return e
}

func (e *MultiFactorRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
