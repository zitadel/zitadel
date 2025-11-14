package schemauser

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventPrefix     = "schemauser."
	CreatedType     = eventPrefix + "created"
	UpdatedType     = eventPrefix + "updated"
	DeletedType     = eventPrefix + "deleted"
	LockedType      = eventPrefix + "locked"
	UnlockedType    = eventPrefix + "unlocked"
	DeactivatedType = eventPrefix + "deactivated"
	ActivatedType   = eventPrefix + "activated"
)

type CreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string          `json:"id"`
	SchemaID              string          `json:"schemaID"`
	SchemaRevision        uint64          `json:"schemaRevision"`
	Data                  json.RawMessage `json:"user,omitempty"`
}

func (e *CreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *CreatedEvent) Payload() interface{} {
	return e
}

func (e *CreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	schemaID string,
	schemaRevision uint64,
	data json.RawMessage,
) *CreatedEvent {
	return &CreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CreatedType,
		),
		SchemaID:       schemaID,
		SchemaRevision: schemaRevision,
		Data:           data,
	}
}

type UpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SchemaID       *string         `json:"schemaID,omitempty"`
	SchemaRevision *uint64         `json:"schemaRevision,omitempty"`
	Data           json.RawMessage `json:"schema,omitempty"`
}

func (e *UpdatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UpdatedEvent) Payload() interface{} {
	return e
}

func (e *UpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
func NewUpdatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []Changes,
) *UpdatedEvent {
	updatedEvent := &UpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UpdatedType,
		),
	}
	for _, change := range changes {
		change(updatedEvent)
	}
	return updatedEvent
}

type Changes func(event *UpdatedEvent)

func ChangeSchemaID(schemaID string) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaID = &schemaID
	}
}
func ChangeSchemaRevision(schemaRevision uint64) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaRevision = &schemaRevision
	}
}

func ChangeData(data json.RawMessage) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.Data = data
	}
}

type DeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *DeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *DeletedEvent) Payload() interface{} {
	return e
}

func (e *DeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DeletedEvent {
	return &DeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DeletedType,
		),
	}
}

type LockedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *LockedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *LockedEvent) Payload() interface{} {
	return e
}

func (e *LockedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewLockedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *LockedEvent {
	return &LockedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LockedType,
		),
	}
}

type UnlockedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *UnlockedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UnlockedEvent) Payload() interface{} {
	return e
}

func (e *UnlockedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUnlockedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *UnlockedEvent {
	return &UnlockedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UnlockedType,
		),
	}
}

type DeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *DeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *DeactivatedEvent) Payload() interface{} {
	return e
}

func (e *DeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DeactivatedEvent {
	return &DeactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DeactivatedType,
		),
	}
}

type ActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *ActivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *ActivatedEvent) Payload() interface{} {
	return e
}

func (e *ActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *ActivatedEvent {
	return &ActivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ActivatedType,
		),
	}
}
