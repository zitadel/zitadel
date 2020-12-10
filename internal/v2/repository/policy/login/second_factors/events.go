package second_factors

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	loginPolicySecondFactorPrefix           = "policy.login.secondfactor."
	LoginPolicySecondFactorAddedEventType   = loginPolicySecondFactorPrefix + "added"
	LoginPolicySecondFactorRemovedEventType = loginPolicySecondFactorPrefix + "removed"
)

type SecondFactorAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MFAType SecondFactorType `json:"mfaType"`
}

func NewSecondFactorAddedEvent(
	base *eventstore.BaseEvent,
	mfaType SecondFactorType,
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

type SecondFactorRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	MFAType              SecondFactorType `json:"mfaType"`
}

func NewSecondFactorRemovedEvent(
	base *eventstore.BaseEvent,
	mfaType SecondFactorType,
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
