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
	emailEventPrefix                 = humanEventPrefix + "email."
	HumanEmailChangedType            = emailEventPrefix + "changed"
	HumanEmailVerifiedType           = emailEventPrefix + "verified"
	HumanEmailVerificationFailedType = emailEventPrefix + "verification.failed"
	HumanEmailCodeAddedType          = emailEventPrefix + "code.added"
	HumanEmailCodeSentType           = emailEventPrefix + "code.sent"
)

type HumanEmailChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	EmailAddress string `json:"email,omitempty"`
}

func (e *HumanEmailChangedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanEmailChangedEvent) Data() interface{} {
	return e
}

func NewHumanHumanEmailChangedEvent(base *eventstore.BaseEvent, emailAddress string) *HumanEmailChangedEvent {
	return &HumanEmailChangedEvent{
		BaseEvent:    *base,
		EmailAddress: emailAddress,
	}
}

func HumanEmailChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailChangedEvent := &HumanEmailChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return emailChangedEvent, nil
}

type HumanEmailVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *HumanEmailVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanEmailVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanEmailVerifiedEvent(base *eventstore.BaseEvent) *HumanEmailVerifiedEvent {
	return &HumanEmailVerifiedEvent{
		BaseEvent: *base,
	}
}

func HumanEmailVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailVerified := &HumanEmailVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsEmailVerified: true,
	}
	return emailVerified, nil
}

type HumanEmailVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanEmailVerificationFailedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanEmailVerificationFailedEvent) Data() interface{} {
	return nil
}

func NewHumanEmailVerificationFailedEvent(base *eventstore.BaseEvent) *HumanEmailVerificationFailedEvent {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *base,
	}
}

func HumanEmailVerificationFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanEmailCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *HumanEmailCodeAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanEmailCodeAddedEvent) Data() interface{} {
	return e
}

func NewHumanEmailCodeAddedEvent(
	base *eventstore.BaseEvent,
	code *crypto.CryptoValue,
	expiry time.Duration) *HumanEmailCodeAddedEvent {
	return &HumanEmailCodeAddedEvent{
		BaseEvent: *base,
		Code:      code,
		Expiry:    expiry,
	}
}

func HumanEmailCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	codeAdded := &HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, codeAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3M0sd", "unable to unmarshal human email code added")
	}

	return codeAdded, nil
}

type HumanEmailCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanEmailCodeSentEvent) CheckPrevious() bool {
	return true
}

func (e *HumanEmailCodeSentEvent) Data() interface{} {
	return nil
}

func NewHumanEmailCodeSentEvent(
	base *eventstore.BaseEvent) *HumanEmailCodeSentEvent {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *base,
	}
}

func HumanEmailCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
