package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func (e *HumanOTPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanOTPAddedEvent) Assets() []*eventstore.Asset {
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
	UserAgentID          string `json:"userAgentID,omitempty"`
}

func (e *HumanOTPVerifiedEvent) Data() interface{} {
	return nil
}

func (e *HumanOTPVerifiedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanOTPVerifiedEvent) Assets() []*eventstore.Asset {
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

func (e *HumanOTPRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanOTPRemovedEvent) Assets() []*eventstore.Asset {
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

func HumanOTPRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *HumanOTPCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanOTPCheckSucceededEvent) Assets() []*eventstore.Asset {
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

func HumanOTPCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check succeeded")
	}
	return otpAdded, nil
}

type HumanOTPCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckFailedEvent) Data() interface{} {
	return e
}

func (e *HumanOTPCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanOTPCheckFailedEvent) Assets() []*eventstore.Asset {
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

func HumanOTPCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	otpAdded := &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, otpAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check failed")
	}
	return otpAdded, nil
}
