package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	passwordEventPrefix             = humanEventPrefix + "password."
	HumanPasswordChangedType        = passwordEventPrefix + "changed"
	HumanPasswordChangeSentType     = passwordEventPrefix + "change.sent"
	HumanPasswordCodeAddedType      = passwordEventPrefix + "code.added"
	HumanPasswordCodeSentType       = passwordEventPrefix + "code.sent"
	HumanPasswordCheckSucceededType = passwordEventPrefix + "check.succeeded"
	HumanPasswordCheckFailedType    = passwordEventPrefix + "check.failed"
)

type HumanPasswordChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired"`
	UserAgentID    string              `json:"userAgentID,omitempty"`
}

func (e *HumanPasswordChangedEvent) Payload() interface{} {
	return e
}

func (e *HumanPasswordChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	secret *crypto.CryptoValue,
	changeRequired bool,
	userAgentID string,
) *HumanPasswordChangedEvent {
	return &HumanPasswordChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordChangedType,
		),
		Secret:         secret,
		ChangeRequired: changeRequired,
		UserAgentID:    userAgentID,
	}
}

func HumanPasswordChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanAdded := &HumanPasswordChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return humanAdded, nil
}

type HumanPasswordCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code             *crypto.CryptoValue     `json:"code,omitempty"`
	Expiry           time.Duration           `json:"expiry,omitempty"`
	NotificationType domain.NotificationType `json:"notificationType,omitempty"`
}

func (e *HumanPasswordCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanPasswordCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	notificationType domain.NotificationType,
) *HumanPasswordCodeAddedEvent {
	return &HumanPasswordCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordCodeAddedType,
		),
		Code:             code,
		Expiry:           expiry,
		NotificationType: notificationType,
	}
}

func HumanPasswordCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanAdded := &HumanPasswordCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ms90d", "unable to unmarshal human password code added")
	}

	return humanAdded, nil
}

type HumanPasswordCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPasswordCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *HumanPasswordCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanPasswordCodeSentEvent {
	return &HumanPasswordCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordCodeSentType,
		),
	}
}

func HumanPasswordCodeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanPasswordCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPasswordChangeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPasswordChangeSentEvent) Payload() interface{} {
	return nil
}

func (e *HumanPasswordChangeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordChangeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanPasswordChangeSentEvent {
	return &HumanPasswordChangeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordChangeSentType,
		),
	}
}

func HumanPasswordChangeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanPasswordChangeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPasswordCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanPasswordCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanPasswordCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanPasswordCheckSucceededEvent {
	return &HumanPasswordCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

func HumanPasswordCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanAdded := &HumanPasswordCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5M9sd", "unable to unmarshal human password check succeeded")
	}

	return humanAdded, nil
}

type HumanPasswordCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanPasswordCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanPasswordCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPasswordCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanPasswordCheckFailedEvent {
	return &HumanPasswordCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}

func HumanPasswordCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanAdded := &HumanPasswordCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4m9fs", "unable to unmarshal human password check failed")
	}

	return humanAdded, nil
}
