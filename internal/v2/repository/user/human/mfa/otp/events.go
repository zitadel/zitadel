package otp

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa"
)

const (
	otpEventPrefix                = eventstore.EventType("user.human.mfa.otp.")
	HumanMFAOTPAddedType          = otpEventPrefix + "added"
	HumanMFAOTPVerifiedType       = otpEventPrefix + "verified"
	HumanMFAOTPRemovedType        = otpEventPrefix + "removed"
	HumanMFAOTPCheckSucceededType = otpEventPrefix + "check.succeeded"
	HumanMFAOTPCheckFailedType    = otpEventPrefix + "check.failed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  mfa.State           `json:"-"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(ctx context.Context,
	secret *crypto.CryptoValue) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPAddedType,
		),
		Secret: secret,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.StateNotReady,
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type VerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	State                mfa.State `json:"-"`
}

func (e *VerifiedEvent) Data() interface{} {
	return nil
}

func NewVerifiedEvent(ctx context.Context) *VerifiedEvent {
	return &VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPVerifiedType,
		),
	}
}

func VerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &VerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.StateReady,
	}, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(ctx context.Context) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPRemovedType,
		),
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type CheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CheckSucceededEvent) Data() interface{} {
	return nil
}

func NewCheckSucceededEvent(ctx context.Context) *CheckSucceededEvent {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckSucceededType,
		),
	}
}

func CheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type CheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CheckFailedEvent) Data() interface{} {
	return nil
}

func NewCheckFailedEvent(ctx context.Context) *CheckFailedEvent {
	return &CheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckFailedType,
		),
	}
}

func CheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
