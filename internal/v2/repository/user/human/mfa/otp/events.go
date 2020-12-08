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

type HumanMFAOTPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  mfa.MFAState        `json:"-"`
}

func (e *HumanMFAOTPAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAOTPAddedEvent) Data() interface{} {
	return e
}

func NewHumanMFAOTPAddedEvent(ctx context.Context,
	secret *crypto.CryptoValue) *HumanMFAOTPAddedEvent {
	return &HumanMFAOTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPAddedType,
		),
		Secret: secret,
	}
}

func HumanMFAOTPAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &HumanMFAOTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.MFAStateNotReady,
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type HumanMFAOTPVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	State                mfa.MFAState `json:"-"`
}

func (e *HumanMFAOTPVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAOTPVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanMFAOTPVerifiedEvent(ctx context.Context) *HumanMFAOTPVerifiedEvent {
	return &HumanMFAOTPVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPVerifiedType,
		),
	}
}

func HumanMFAOTPVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.MFAStateReady,
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

func NewHumanMFAOTPRemovedEvent(ctx context.Context) *HumanMFAOTPRemovedEvent {
	return &HumanMFAOTPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPRemovedType,
		),
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

func NewHumanMFAOTPCheckSucceededEvent(ctx context.Context) *HumanMFAOTPCheckSucceededEvent {
	return &HumanMFAOTPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckSucceededType,
		),
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

func NewHumanMFAOTPCheckFailedEvent(ctx context.Context) *HumanMFAOTPCheckFailedEvent {
	return &HumanMFAOTPCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAOTPCheckFailedType,
		),
	}
}

func HumanMFAOTPCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
