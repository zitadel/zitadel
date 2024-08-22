package schemauser

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventPrefix = "user_schema.user."
	CreatedType = eventPrefix + "created"
	UpdatedType = eventPrefix + "updated"
	DeletedType = eventPrefix + "deleted"
)

type CreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string          `json:"id"`
	SchemaType            string          `json:"schemaType"`
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

/*
func (e *CreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUserIDUniqueConstraint(e.SchemaID)}
}
*/

func NewCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	schemaType string,
	schemaRevision uint64,
	data json.RawMessage,
) *CreatedEvent {
	return &CreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CreatedType,
		),
		SchemaType:     schemaType,
		SchemaRevision: schemaRevision,
		Data:           data,
	}
}

type UpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SchemaType     *string         `json:"schemaType,omitempty"`
	SchemaRevision *uint64         `json:"schemaRevision,omitempty"`
	Data           json.RawMessage `json:"schema,omitempty"`
	oldSchemaType  string
	oldRevision    uint64
}

func (e *UpdatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UpdatedEvent) Payload() interface{} {
	return e
}

/*
func (e *UpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldSchemaType == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveUserIDUniqueConstraint(e.Agg.ID),
		NewAddUserIDUniqueConstraint(*e.Agg.ID),
	}
}*/

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

func ChangeSchemaType(oldSchemaType, schemaType string) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaType = &schemaType
		e.oldSchemaType = oldSchemaType
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
