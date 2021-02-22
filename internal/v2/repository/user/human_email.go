package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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

func (e *HumanEmailChangedEvent) Data() interface{} {
	return e
}

func (e *HumanEmailChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanEmailChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailChangedEvent {
	return &HumanEmailChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailChangedType,
		),
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

func (e *HumanEmailVerifiedEvent) Data() interface{} {
	return nil
}

func (e *HumanEmailVerifiedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanEmailVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailVerifiedEvent {
	return &HumanEmailVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanEmailVerificationFailedEvent) Data() interface{} {
	return nil
}

func (e *HumanEmailVerificationFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanEmailVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailVerificationFailedEvent {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanEmailCodeAddedEvent) Data() interface{} {
	return e
}

func (e *HumanEmailCodeAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanEmailCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration) *HumanEmailCodeAddedEvent {
	return &HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanEmailCodeSentEvent) Data() interface{} {
	return nil
}

func (e *HumanEmailCodeSentEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanEmailCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailCodeSentEvent {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailCodeSentType,
		),
	}
}

func HumanEmailCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
