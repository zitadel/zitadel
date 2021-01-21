package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
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

func SecondFactorAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &SecondFactorAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-Lp0dE", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *SecondFactorAddedEvent) Data() interface{} {
	return e
}

func (e *SecondFactorAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func SecondFactorRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &SecondFactorRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-5M9gd", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *SecondFactorRemovedEvent) Data() interface{} {
	return e
}

func (e *SecondFactorRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func MultiFactorAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MultiFactorAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-5Ms90", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *MultiFactorAddedEvent) Data() interface{} {
	return e
}

func (e *MultiFactorAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func MultiFactorRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MultiFactorRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-1N8sd", "unable to unmarshal policy")
	}

	return e, nil
}

func (e *MultiFactorRemovedEvent) Data() interface{} {
	return e
}

func (e *MultiFactorRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}
