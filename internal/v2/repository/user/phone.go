package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	phoneEventPrefix                 = humanEventPrefix + "phone."
	HumanPhoneChangedType            = phoneEventPrefix + "changed"
	HumanPhoneRemovedType            = phoneEventPrefix + "removed"
	HumanPhoneVerifiedType           = phoneEventPrefix + "verified"
	HumanPhoneVerificationFailedType = phoneEventPrefix + "verification.failed"
	HumanPhoneCodeAddedType          = phoneEventPrefix + "code.added"
	HumanPhoneCodeSentType           = phoneEventPrefix + "code.sent"
)

type HumanPhoneChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PhoneNumber string `json:"phone,omitempty"`
}

func (e *HumanPhoneChangedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPhoneChangedEvent) Data() interface{} {
	return e
}

func NewHumanPhoneChangedEvent(base *eventstore.BaseEvent, phone string) *HumanPhoneChangedEvent {
	return &HumanPhoneChangedEvent{
		BaseEvent:   *base,
		PhoneNumber: phone,
	}
}

func HumanPhoneChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	phoneChangedEvent := &HumanPhoneChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, phoneChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human phone changed")
	}

	return phoneChangedEvent, nil
}

type HumanPhoneRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPhoneRemovedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPhoneRemovedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneRemovedEvent(base *eventstore.BaseEvent) *HumanPhoneRemovedEvent {
	return &HumanPhoneRemovedEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneRemovedEventMapper(event *repository.Event) eventstore.EventReader {
	return &HumanPhoneChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
}

type HumanPhoneVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *HumanPhoneVerifiedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPhoneVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneVerifiedEvent(base *eventstore.BaseEvent) *HumanPhoneVerifiedEvent {
	return &HumanPhoneVerifiedEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneVerifiedEventMapper(event *repository.Event) eventstore.EventReader {
	return &HumanPhoneVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsPhoneVerified: true,
	}
}
