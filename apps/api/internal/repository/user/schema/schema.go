package schema

import (
	"context"
	"encoding/json"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventPrefix     = "user_schema."
	CreatedType     = eventPrefix + "created"
	UpdatedType     = eventPrefix + "updated"
	DeactivatedType = eventPrefix + "deactivated"
	ReactivatedType = eventPrefix + "reactivated"
	DeletedType     = eventPrefix + "deleted"

	uniqueSchemaType = "user_schema_type"
)

func NewAddSchemaTypeUniqueConstraint(schemaType string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueSchemaType,
		schemaType,
		"Errors.UserSchema.Type.AlreadyExists")
}

func NewRemoveSchemaTypeUniqueConstraint(schemaType string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueSchemaType,
		schemaType,
	)
}

type CreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SchemaType             string                     `json:"schemaType"`
	Schema                 json.RawMessage            `json:"schema,omitempty"`
	PossibleAuthenticators []domain.AuthenticatorType `json:"possibleAuthenticators,omitempty"`
}

func (e *CreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *CreatedEvent) Payload() interface{} {
	return e
}

func (e *CreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddSchemaTypeUniqueConstraint(e.SchemaType)}
}

func NewCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	schemaType string,
	schema json.RawMessage,
	possibleAuthenticators []domain.AuthenticatorType,
) *CreatedEvent {
	return &CreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CreatedType,
		),
		SchemaType:             schemaType,
		Schema:                 schema,
		PossibleAuthenticators: possibleAuthenticators,
	}
}

type UpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SchemaType             *string                    `json:"schemaType,omitempty"`
	Schema                 json.RawMessage            `json:"schema,omitempty"`
	PossibleAuthenticators []domain.AuthenticatorType `json:"possibleAuthenticators,omitempty"`
	SchemaRevision         *uint64                    `json:"schemaRevision,omitempty"`
	oldSchemaType          string
	oldRevision            uint64
}

func (e *UpdatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UpdatedEvent) Payload() interface{} {
	return e
}

func (e *UpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldSchemaType == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveSchemaTypeUniqueConstraint(e.oldSchemaType),
		NewAddSchemaTypeUniqueConstraint(*e.SchemaType),
	}
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

func ChangeSchema(schema json.RawMessage) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.Schema = schema
	}
}

func ChangePossibleAuthenticators(possibleAuthenticators []domain.AuthenticatorType) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.PossibleAuthenticators = possibleAuthenticators
	}
}

func IncreaseRevision(oldRevision uint64) func(event *UpdatedEvent) {
	return func(e *UpdatedEvent) {
		e.SchemaRevision = gu.Ptr(oldRevision + 1)
		e.oldRevision = oldRevision
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

type ReactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *ReactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *ReactivatedEvent) Payload() interface{} {
	return e
}

func (e *ReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewReactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *ReactivatedEvent {
	return &ReactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ReactivatedType,
		),
	}
}

type DeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	schemaType string
}

func (e *DeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *DeletedEvent) Payload() interface{} {
	return e
}

func (e *DeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveSchemaTypeUniqueConstraint(e.schemaType),
	}
}

func NewDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	schemaType string,
) *DeletedEvent {
	return &DeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DeletedType,
		),
		schemaType: schemaType,
	}
}
