package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
)

const (
	passwordEventPrefix             = humanEventPrefix + "password."
	HumanPasswordChangedType        = passwordEventPrefix + "changed"
	HumanPasswordCodeAddedType      = passwordEventPrefix + "code.added"
	HumanPasswordCodeSentType       = passwordEventPrefix + "code.sent"
	HumanPasswordCheckSucceededType = passwordEventPrefix + "check.succeeded"
	HumanPasswordCheckFailedType    = passwordEventPrefix + "check.failed"
)

type HumanPasswordChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

func (e *HumanPasswordChangedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPasswordChangedEvent) Data() interface{} {
	return e
}

func HumanPasswordChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &HumanPasswordChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return humanAdded, nil
}

type HumanPasswordCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code             *crypto.CryptoValue `json:"code,omitempty"`
	Expiry           time.Duration       `json:"expiry,omitempty"`
	NotificationType int32               `json:"notificationType,omitempty"`
}

func (e *HumanPasswordCodeAddedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPasswordCodeAddedEvent) Data() interface{} {
	return e
}

func HumanPasswordCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &HumanPasswordCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Ms90d", "unable to unmarshal human password code added")
	}

	return humanAdded, nil
}

type HumanPasswordCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPasswordCodeSentEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPasswordCodeSentEvent) Data() interface{} {
	return nil
}

func HumanPasswordCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPasswordCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPasswordCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPasswordCheckSucceededEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPasswordCheckSucceededEvent) Data() interface{} {
	return nil
}

func HumanPasswordCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPasswordCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPasswordCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPasswordCheckFailedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanPasswordCheckFailedEvent) Data() interface{} {
	return nil
}

func HumanPasswordCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanPasswordCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
