package phone

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
)

const (
	phoneEventPrefix                 = eventstore.EventType("user.human.phone.")
	HumanPhoneChangedType            = phoneEventPrefix + "changed"
	HumanPhoneRemovedType            = phoneEventPrefix + "removed"
	HumanPhoneVerifiedType           = phoneEventPrefix + "verified"
	HumanPhoneVerificationFailedType = phoneEventPrefix + "verification.failed"
	HumanPhoneCodeAddedType          = phoneEventPrefix + "code.added"
	HumanPhoneCodeSentType           = phoneEventPrefix + "code.sent"
)

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PhoneNumber string `json:"phone,omitempty"`
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewHumanPhoneChangedEvent(ctx context.Context, phone string) *ChangedEvent {
	return &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneChangedType,
		),
		PhoneNumber: phone,
	}
}

func HumanPhoneChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	phoneChangedEvent := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, phoneChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human phone changed")
	}

	return phoneChangedEvent, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneRemovedEvent(ctx context.Context) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneRemovedType,
		),
	}
}

func HumanPhoneRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type VerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *VerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *VerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneVerifiedEvent(ctx context.Context) *VerifiedEvent {
	return &VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneVerifiedType,
		),
	}
}

func HumanPhoneVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &VerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsPhoneVerified: true,
	}, nil
}

type VerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *VerificationFailedEvent) CheckPrevious() bool {
	return true
}

func (e *VerificationFailedEvent) Data() interface{} {
	return nil
}

func NewHumanPhoneVerificationFailedEvent(ctx context.Context) *VerificationFailedEvent {
	return &VerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneVerificationFailedType,
		),
	}
}

func HumanPhoneVerificationFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &VerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type CodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *CodeAddedEvent) CheckPrevious() bool {
	return true
}

func (e *CodeAddedEvent) Data() interface{} {
	return e
}

func NewHumanPhoneCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *CodeAddedEvent {
	return &CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func HumanPhoneCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	codeAdded := &CodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, codeAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-6Ms9d", "unable to unmarshal human phone code added")
	}

	return codeAdded, nil
}

type CodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CodeSentEvent) CheckPrevious() bool {
	return false
}

func (e *CodeSentEvent) Data() interface{} {
	return e
}

func NewHumanPhoneCodeSentEvent(ctx context.Context) *CodeSentEvent {
	return &CodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPhoneCodeSentType,
		),
	}
}

func HumanPhoneCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
