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

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	EmailAddress string `json:"email,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	ctx context.Context,
	current *WriteModel,
	emailAddress string,
) *ChangedEvent {
	e := &ChangedEvent{
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

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailChangedEvent := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return emailChangedEvent, nil
}

type VerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *VerifiedEvent) Data() interface{} {
	return nil
}

func NewVerifiedEvent(ctx context.Context) *VerifiedEvent {
	return &VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailVerifiedType,
		),
	}
}

func VerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailVerified := &VerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsEmailVerified: true,
	}
	return emailVerified, nil
}

type VerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *VerificationFailedEvent) Data() interface{} {
	return nil
}

func NewVerificationFailedEvent(ctx context.Context) *VerificationFailedEvent {
	return &VerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailVerificationFailedType,
		),
	}
}

func VerificationFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &VerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type CodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *CodeAddedEvent) Data() interface{} {
	return e
}

func NewCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration) *CodeAddedEvent {
	return &CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func CodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	codeAdded := &CodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, codeAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3M0sd", "unable to unmarshal human email code added")
	}

	return codeAdded, nil
}

type CodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CodeSentEvent) Data() interface{} {
	return nil
}

func NewCodeSentEvent(ctx context.Context) *CodeSentEvent {
	return &CodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanEmailCodeSentType,
		),
	}
}

func CodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
