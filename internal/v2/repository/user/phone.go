package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
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
	return true
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
	return true
}

func (e *HumanPhoneRemovedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneRemovedEvent(base *eventstore.BaseEvent) *HumanPhoneRemovedEvent {
	return &HumanPhoneRemovedEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPhoneChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPhoneVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *HumanPhoneVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanPhoneVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneVerifiedEvent(base *eventstore.BaseEvent) *HumanPhoneVerifiedEvent {
	return &HumanPhoneVerifiedEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPhoneVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsPhoneVerified: true,
	}, nil
}

type HumanPhoneVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPhoneVerificationFailedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanPhoneVerificationFailedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneVerificationFailedEvent(base *eventstore.BaseEvent) *HumanPhoneVerificationFailedEvent {
	return &HumanPhoneVerificationFailedEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneVerificationFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPhoneVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPhoneCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *HumanPhoneCodeAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanPhoneCodeAddedEvent) Data() interface{} {
	return e
}

func NewHumanPhoneCodeAddedEvent(
	base *eventstore.BaseEvent,
	code *crypto.CryptoValue,
	expiry time.Duration) *HumanPhoneCodeAddedEvent {
	return &HumanPhoneCodeAddedEvent{
		BaseEvent: *base,
		Code:      code,
		Expiry:    expiry,
	}
}

func HumanPhoneCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	codeAdded := &HumanPhoneCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, codeAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-6Ms9d", "unable to unmarshal human phone code added")
	}

	return codeAdded, nil
}

type HumanPhoneCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPhoneCodeSentEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPhoneCodeSentEvent) Data() interface{} {
	return e
}

func NewHumanPhoneCodeSentEvent(
	base *eventstore.BaseEvent) *HumanPhoneCodeSentEvent {
	return &HumanPhoneCodeSentEvent{
		BaseEvent: *base,
	}
}

func HumanPhoneCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPhoneCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
