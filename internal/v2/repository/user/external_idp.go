package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	externalIDPEventPrefix   = humanEventPrefix + "externalidp."
	externalLoginEventPrefix = humanEventPrefix + "externallogin."

	HumanExternalIDPReservedType = externalIDPEventPrefix + "reserved"
	HumanExternalIDPReleasedType = externalIDPEventPrefix + "released"

	HumanExternalIDPAddedType          = externalIDPEventPrefix + "added"
	HumanExternalIDPRemovedType        = externalIDPEventPrefix + "removed"
	HumanExternalIDPCascadeRemovedType = externalIDPEventPrefix + "cascade.removed"

	HumanExternalLoginCheckSucceededType = externalLoginEventPrefix + "check.succeeded"
)

type HumanExternalIDPReservedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanExternalIDPReservedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPReservedEvent) Data() interface{} {
	return nil
}

func NewHumanExternalIDPReservedEvent(base *eventstore.BaseEvent) *HumanExternalIDPReservedEvent {
	return &HumanExternalIDPReservedEvent{
		BaseEvent: *base,
	}
}

type HumanExternalIDPReleasedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanExternalIDPReleasedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPReleasedEvent) Data() interface{} {
	return nil
}

func NewHumanExternalIDPReleasedEvent(base *eventstore.BaseEvent) *HumanExternalIDPReleasedEvent {
	return &HumanExternalIDPReleasedEvent{
		BaseEvent: *base,
	}
}

type HumanExternalIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func (e *HumanExternalIDPAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPAddedEvent) Data() interface{} {
	return e
}

func NewHumanExternalIDPAddedEvent(base *eventstore.BaseEvent, idpConfigID, displayName string) *HumanExternalIDPAddedEvent {
	return &HumanExternalIDPAddedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
		DisplayName: displayName,
	}
}

func HumanExternalIDPAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-6M9sd", "unable to unmarshal user external idp added")
	}

	return e, nil
}

type HumanExternalIDPRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *HumanExternalIDPRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPRemovedEvent) Data() interface{} {
	return e
}

func NewHumanExternalIDPRemovedEvent(base *eventstore.BaseEvent, idpConfigID string) *HumanExternalIDPRemovedEvent {
	return &HumanExternalIDPRemovedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func HumanExternalIDPRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal user external idp removed")
	}

	return e, nil
}

type HumanExternalIDPCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *HumanExternalIDPCascadeRemovedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanExternalIDPCascadeRemovedEvent) Data() interface{} {
	return e
}

func NewHumanExternalIDPCascadeRemovedEvent(base *eventstore.BaseEvent, idpConfigID string) *HumanExternalIDPCascadeRemovedEvent {
	return &HumanExternalIDPCascadeRemovedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func HumanExternalIDPCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0sd", "unable to unmarshal user external idp cascade removed")
	}

	return e, nil
}

type HumanExternalLoginCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanExternalLoginCheckSucceededEvent) CheckPrevious() bool {
	return false
}

func (e *HumanExternalLoginCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewHumanExternalLoginCheckSucceededEvent(base *eventstore.BaseEvent) *HumanExternalLoginCheckSucceededEvent {
	return &HumanExternalLoginCheckSucceededEvent{
		BaseEvent: *base,
	}
}

func HumanExternalLoginCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanExternalLoginCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
