package schemauser

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventPrefix = "user."
	CreatedType = eventPrefix + "created"
	UpdatedType = eventPrefix + "updated"
	DeletedType = eventPrefix + "deleted"
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
	oldSchemaID    string
	oldRevision    uint64
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

func ChangeSchemaID(oldSchemaID, schemaID string) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaID = &schemaID
		e.oldSchemaID = oldSchemaID
	}
}
func ChangeSchemaRevision(oldSchemaRevision, schemaRevision uint64) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaRevision = &schemaRevision
		e.oldRevision = oldSchemaRevision
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
