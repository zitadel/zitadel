package email

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
	emailEventPrefix                 = eventstore.EventType("user.human.email.")
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

func NewHumanEmailChangedEvent(
	ctx context.Context,
	current *HumanEmailWriteModel,
	emailAddress string,
) *HumanEmailChangedEvent {
	e := &HumanEmailChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailChangedType,
		),
	}
	if current.Email != emailAddress {
		e.EmailAddress = emailAddress
	}
	return e
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

func NewHumanEmailVerifiedEvent(ctx context.Context) *HumanEmailVerifiedEvent {
	return &HumanEmailVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailVerifiedType,
		),
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

func NewHumanEmailVerificationFailedEvent(ctx context.Context) *HumanEmailVerificationFailedEvent {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailVerificationFailedType,
		),
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
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration) *HumanEmailCodeAddedEvent {
	return &HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
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

func NewHumanEmailCodeSentEvent(ctx context.Context) *HumanEmailCodeSentEvent {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailCodeSentType,
		),
	}
}

func HumanEmailCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
