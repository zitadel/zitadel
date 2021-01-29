package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	uniqueExternalIDPType    = "external_idps"
	externalIDPEventPrefix   = humanEventPrefix + "externalidp."
	externalLoginEventPrefix = humanEventPrefix + "externallogin."

	//TODO: Handle unique Aggregate
	HumanExternalIDPReservedType = externalIDPEventPrefix + "reserved"
	HumanExternalIDPReleasedType = externalIDPEventPrefix + "released"

	HumanExternalIDPAddedType          = externalIDPEventPrefix + "added"
	HumanExternalIDPRemovedType        = externalIDPEventPrefix + "removed"
	HumanExternalIDPCascadeRemovedType = externalIDPEventPrefix + "cascade.removed"

	HumanExternalLoginCheckSucceededType = externalLoginEventPrefix + "check.succeeded"
)

func NewAddExternalIDPUniqueConstraint(idpConfigID, externalUserID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueExternalIDPType,
		idpConfigID+externalUserID,
		"Errors.User.ExternalIDP.AlreadyExists")
}

func NewRemoveExternalIDPUniqueConstraint(idpConfigID, externalUserID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueExternalIDPType,
		idpConfigID+externalUserID)
}

type HumanExternalIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId,omitempty"`
	UserID      string `json:"userId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func (e *HumanExternalIDPAddedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddExternalIDPUniqueConstraint(e.IDPConfigID, e.UserID)}
}

func NewHumanExternalIDPAddedEvent(ctx context.Context, idpConfigID, displayName string) *HumanExternalIDPAddedEvent {
	return &HumanExternalIDPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPAddedType,
		),
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
	UserID      string `json:"userId,omitempty"`
}

func (e *HumanExternalIDPRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveExternalIDPUniqueConstraint(e.IDPConfigID, e.UserID)}
}

func NewHumanExternalIDPRemovedEvent(ctx context.Context, idpConfigID, externalUserID string) *HumanExternalIDPRemovedEvent {
	return &HumanExternalIDPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPRemovedType,
		),
		IDPConfigID: idpConfigID,
		UserID:      externalUserID,
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
	UserID      string `json:"userId,omitempty"`
}

func (e *HumanExternalIDPCascadeRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanExternalIDPCascadeRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveExternalIDPUniqueConstraint(e.IDPConfigID, e.UserID)}
}

func NewHumanExternalIDPCascadeRemovedEvent(ctx context.Context, idpConfigID, externalUserID string) *HumanExternalIDPCascadeRemovedEvent {
	return &HumanExternalIDPCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPCascadeRemovedType,
		),
		IDPConfigID: idpConfigID,
		UserID:      externalUserID,
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

func NewHumanExternalIDPCheckSucceededEvent(ctx context.Context, info *AuthRequestInfo) *HumanExternalIDPCheckSucceededEvent {
	return &HumanExternalIDPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
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
