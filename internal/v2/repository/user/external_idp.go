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

type HumanExternalIDPAdded struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func (e *HumanExternalIDPAdded) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPAdded) Data() interface{} {
	return e
}

func NewHumanExternalIDPAddedEvent(base *eventstore.BaseEvent, idpConfigID, displayName string) *HumanExternalIDPAdded {
	return &HumanExternalIDPAdded{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
		DisplayName: displayName,
	}
}

func HumanExternalIDPAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPAdded{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-6M9sd", "unable to unmarshal user external idp added")
	}

	return e, nil
}

type HumanExternalIDPRemoved struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *HumanExternalIDPRemoved) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPRemoved) Data() interface{} {
	return e
}

func NewHumanExternalIDPRemovedEvent(base *eventstore.BaseEvent, idpConfigID string) *HumanExternalIDPRemoved {
	return &HumanExternalIDPRemoved{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func HumanExternalIDPRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPRemoved{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal user external idp removed")
	}

	return e, nil
}

type HumanExternalIDPCascadeRemoved struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *HumanExternalIDPCascadeRemoved) CheckPrevious() bool {
	return true
}

func (e *HumanExternalIDPCascadeRemoved) Data() interface{} {
	return e
}

func NewHumanExternalIDPCascadeRemovedEvent(base *eventstore.BaseEvent, idpConfigID string) *HumanExternalIDPCascadeRemoved {
	return &HumanExternalIDPCascadeRemoved{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func HumanExternalIDPCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPCascadeRemoved{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0sd", "unable to unmarshal user external idp cascade removed")
	}

	return e, nil
}
