package multi_factors

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	loginPolicyMultiFactorPrefix           = "policy.login.multifactor."
	LoginPolicyMultiFactorAddedEventType   = loginPolicyMultiFactorPrefix + "added"
	LoginPolicyMultiFactorRemovedEventType = loginPolicyMultiFactorPrefix + "removed"
)

type MultiFactorAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MFAType MultiFactorType `json:"mfaType"`
}

func NewMultiFactorAddedEvent(
	base *eventstore.BaseEvent,
	mfaType MultiFactorType,
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

func (e *MultiFactorAddedEvent) CheckPrevious() bool {
	return true
}

func (e *MultiFactorAddedEvent) Data() interface{} {
	return e
}

type MultiFactorRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	MFAType              MultiFactorType `json:"mfaType"`
}

func NewMultiFactorRemovedEvent(
	base *eventstore.BaseEvent,
	mfaType MultiFactorType,
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

func (e *MultiFactorRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *MultiFactorRemovedEvent) Data() interface{} {
	return e
}
