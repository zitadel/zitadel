package password

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/user/human"
	"time"
)

const (
	passwordEventPrefix             = eventstore.EventType("user.human.password.")
	HumanPasswordChangedType        = passwordEventPrefix + "changed"
	HumanPasswordCodeAddedType      = passwordEventPrefix + "code.added"
	HumanPasswordCodeSentType       = passwordEventPrefix + "code.sent"
	HumanPasswordCheckSucceededType = passwordEventPrefix + "check.succeeded"
	HumanPasswordCheckFailedType    = passwordEventPrefix + "check.failed"
)

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	ctx context.Context,
	secret *crypto.CryptoValue,
	changeRequired bool,
) *ChangedEvent {
	return &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordChangedType,
		),
		Secret:         secret,
		ChangeRequired: changeRequired,
	}
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return humanAdded, nil
}

type CodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code             *crypto.CryptoValue    `json:"code,omitempty"`
	Expiry           time.Duration          `json:"expiry,omitempty"`
	NotificationType human.NotificationType `json:"notificationType,omitempty"`
}

func (e *CodeAddedEvent) Data() interface{} {
	return e
}

func NewPasswordCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
	notificationType human.NotificationType,
) *CodeAddedEvent {
	return &CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordCodeAddedType,
		),
		Code:             code,
		Expiry:           expiry,
		NotificationType: notificationType,
	}
}

func CodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &CodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ms90d", "unable to unmarshal human password code added")
	}

	return humanAdded, nil
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
			HumanPasswordCodeSentType,
		),
	}
}

func CodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CodeSentEvent{
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
			HumanPasswordCheckSucceededType,
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
			HumanPasswordCheckFailedType,
		),
	}
}

func CheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
