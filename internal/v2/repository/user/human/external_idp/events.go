package external_idp

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	externalIDPEventPrefix   = eventstore.EventType("user.human.externalidp.")
	externalLoginEventPrefix = eventstore.EventType("user.human.externallogin.")

	//TODO: Handle unique Aggregate
	HumanExternalIDPReservedType = externalIDPEventPrefix + "reserved"
	HumanExternalIDPReleasedType = externalIDPEventPrefix + "released"

	HumanExternalIDPAddedType          = externalIDPEventPrefix + "added"
	HumanExternalIDPRemovedType        = externalIDPEventPrefix + "removed"
	HumanExternalIDPCascadeRemovedType = externalIDPEventPrefix + "cascade.removed"

	HumanExternalLoginCheckSucceededType = externalLoginEventPrefix + "check.succeeded"
)

type ReservedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ReservedEvent) CheckPrevious() bool {
	return true
}

func (e *ReservedEvent) Data() interface{} {
	return nil
}

func NewReservedEvent(ctx context.Context) *ReservedEvent {
	return &ReservedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPReservedType,
		),
	}
}

type ReleasedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ReleasedEvent) CheckPrevious() bool {
	return true
}

func (e *ReleasedEvent) Data() interface{} {
	return nil
}

func NewReleasedEvent(ctx context.Context) *ReleasedEvent {
	return &ReleasedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPReleasedType,
		),
	}
}

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(ctx context.Context, idpConfigID, displayName string) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPAddedType,
		),
		IDPConfigID: idpConfigID,
		DisplayName: displayName,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-6M9sd", "unable to unmarshal user external idp added")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewRemovedEvent(ctx context.Context, idpConfigID string) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPRemovedType,
		),
		IDPConfigID: idpConfigID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal user external idp removed")
	}

	return e, nil
}

type CascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *CascadeRemovedEvent) CheckPrevious() bool {
	return false
}

func (e *CascadeRemovedEvent) Data() interface{} {
	return e
}

func NewCascadeRemovedEvent(ctx context.Context, idpConfigID string) *CascadeRemovedEvent {
	return &CascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalIDPCascadeRemovedType,
		),
		IDPConfigID: idpConfigID,
	}
}

func CascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &CascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0sd", "unable to unmarshal user external idp cascade removed")
	}

	return e, nil
}

type CheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CheckSucceededEvent) CheckPrevious() bool {
	return false
}

func (e *CheckSucceededEvent) Data() interface{} {
	return nil
}

func NewCheckSucceededEvent(ctx context.Context) *CheckSucceededEvent {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanExternalLoginCheckSucceededType,
		),
	}
}

func CheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
