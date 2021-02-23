package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	UniqueExternalIDPType    = "external_idps"
	externalIDPEventPrefix   = humanEventPrefix + "externalidp."
	externalLoginEventPrefix = humanEventPrefix + "externallogin."

	HumanExternalIDPAddedType          = externalIDPEventPrefix + "added"
	HumanExternalIDPRemovedType        = externalIDPEventPrefix + "removed"
	HumanExternalIDPCascadeRemovedType = externalIDPEventPrefix + "cascade.removed"

	HumanExternalLoginCheckSucceededType = externalLoginEventPrefix + "check.succeeded"
)

func NewAddExternalIDPUniqueConstraint(idpConfigID, externalUserID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueExternalIDPType,
		idpConfigID+externalUserID,
		"Errors.User.ExternalIDP.AlreadyExists")
}

func NewRemoveExternalIDPUniqueConstraint(idpConfigID, externalUserID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueExternalIDPType,
		idpConfigID+externalUserID)
}

type HumanExternalIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID    string `json:"idpConfigId,omitempty"`
	ExternalUserID string `json:"userId,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
}

func (e *HumanExternalIDPAddedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddExternalIDPUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewHumanExternalIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	displayName,
	externalUserID string,
) *HumanExternalIDPAddedEvent {
	return &HumanExternalIDPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanExternalIDPAddedType,
		),
		IDPConfigID:    idpConfigID,
		DisplayName:    displayName,
		ExternalUserID: externalUserID,
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

	IDPConfigID    string `json:"idpConfigId"`
	ExternalUserID string `json:"userId,omitempty"`
}

func (e *HumanExternalIDPRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveExternalIDPUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewHumanExternalIDPRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	externalUserID string,
) *HumanExternalIDPRemovedEvent {
	return &HumanExternalIDPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanExternalIDPRemovedType,
		),
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
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

	IDPConfigID    string `json:"idpConfigId"`
	ExternalUserID string `json:"userId,omitempty"`
}

func (e *HumanExternalIDPCascadeRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPCascadeRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveExternalIDPUniqueConstraint(e.IDPConfigID, e.ExternalUserID)}
}

func NewHumanExternalIDPCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	externalUserID string,
) *HumanExternalIDPCascadeRemovedEvent {
	return &HumanExternalIDPCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanExternalIDPCascadeRemovedType,
		),
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
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

type HumanExternalIDPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanExternalIDPCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanExternalIDPCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *HumanExternalIDPCheckSucceededEvent {
	return &HumanExternalIDPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanExternalLoginCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

func HumanExternalIDPCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &HumanExternalIDPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0sd", "unable to unmarshal user external idp check succeeded")
	}

	return e, nil
}
