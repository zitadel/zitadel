package user

import (
	"context"
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

type HumanOTPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
}

func (e *HumanOTPAddedEvent) Data() interface{} {
	return e
}

func NewHumanOTPAddedEvent(ctx context.Context,
	secret *crypto.CryptoValue) *HumanOTPAddedEvent {
	return &HumanOTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPAddedType,
		),
		Secret: secret,
	}
}

func HumanOTPAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &HumanOTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type HumanOTPVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanOTPVerifiedEvent(ctx context.Context) *HumanOTPVerifiedEvent {
	return &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPVerifiedType,
		),
	}
}

func HumanOTPVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPRemovedEvent) Data() interface{} {
	return nil
}

func NewHumanOTPRemovedEvent(ctx context.Context) *HumanOTPRemovedEvent {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPRemovedType,
		),
	}
}

func HumanOTPRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewHumanOTPCheckSucceededEvent(ctx context.Context) *HumanOTPCheckSucceededEvent {
	return &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckSucceededType,
		),
	}
}

func HumanOTPCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPCheckFailedEvent) Data() interface{} {
	return nil
}

func NewHumanOTPCheckFailedEvent(ctx context.Context) *HumanOTPCheckFailedEvent {
	return &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckFailedType,
		),
	}
}

func HumanOTPCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
