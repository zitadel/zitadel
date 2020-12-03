package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	otpEventPrefix                = mfaEventPrefix + "otp."
	HumanMFAOTPAddedType          = otpEventPrefix + "added"
	HumanMFAOTPVerifiedType       = otpEventPrefix + "verified"
	HumanMFAOTPRemovedType        = otpEventPrefix + "removed"
	HumanMFAOTPCheckSucceededType = otpEventPrefix + "check.succeeded"
	HumanMFAOTPCheckFailedType    = otpEventPrefix + "check.failed"
)

type HumanMFAOTPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  MFAState            `json:"-"`
}

func (e *HumanMFAOTPAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAOTPAddedEvent) Data() interface{} {
	return e
}

func NewHumanMFAOTPAddedEvent(base *eventstore.BaseEvent,
	secret *crypto.CryptoValue) *HumanMFAOTPAddedEvent {
	return &HumanMFAOTPAddedEvent{
		BaseEvent: *base,
		Secret:    secret,
	}
}

func HumanMFAOTPAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &HumanMFAOTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     MFAStateNotReady,
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type HumanMFAOTPVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	State                MFAState `json:"-"`
}

func (e *HumanMFAOTPVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAOTPVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanMFAOTPVerifiedEvent(base *eventstore.BaseEvent) *HumanMFAOTPVerifiedEvent {
	return &HumanMFAOTPVerifiedEvent{
		BaseEvent: *base,
	}
}

func HumanMFAOTPVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     MFAStateReady,
	}, nil
}

type HumanMFAOTPRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanMFAOTPRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAOTPRemovedEvent) Data() interface{} {
	return nil
}

func NewHumanMFAOTPRemovedEvent(base *eventstore.BaseEvent) *HumanMFAOTPRemovedEvent {
	return &HumanMFAOTPRemovedEvent{
		BaseEvent: *base,
	}
}

func HumanMFAOTPRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanMFAOTPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanMFAOTPCheckSucceededEvent) CheckPrevious() bool {
	return false
}

func (e *HumanMFAOTPCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewHumanMFAOTPCheckSucceededEvent(base *eventstore.BaseEvent) *HumanMFAOTPCheckSucceededEvent {
	return &HumanMFAOTPCheckSucceededEvent{
		BaseEvent: *base,
	}
}

func HumanMFAOTPCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanMFAOTPCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanMFAOTPCheckFailedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanMFAOTPCheckFailedEvent) Data() interface{} {
	return nil
}

func NewHumanMFAOTPCheckFailedEvent(base *eventstore.BaseEvent) *HumanMFAOTPCheckFailedEvent {
	return &HumanMFAOTPCheckFailedEvent{
		BaseEvent: *base,
	}
}

func HumanMFAOTPCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
