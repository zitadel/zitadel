package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
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

func (e *HumanOTPAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	secret *crypto.CryptoValue,
) *HumanOTPAddedEvent {
	return &HumanOTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPAddedType,
		),
		Secret: secret,
	}
}

func HumanOTPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type HumanOTPVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserAgentID          string `json:"userAgentID,omitempty"`
}

func (e *HumanOTPVerifiedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPVerifiedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userAgentID string,
) *HumanOTPVerifiedEvent {
	return &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPVerifiedType,
		),
		UserAgentID: userAgentID,
	}
}

func HumanOTPVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPRemovedEvent {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPRemovedType,
		),
	}
}

func HumanOTPRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPCheckSucceededEvent {
	return &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

func HumanOTPCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check succeeded")
	}
	return otpAdded, nil
}

type HumanOTPCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPCheckFailedEvent {
	return &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}

func HumanOTPCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check failed")
	}
	return otpAdded, nil
}
