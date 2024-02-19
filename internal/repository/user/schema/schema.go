package schema

import (
	"context"

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

	uniqueSchemaType = "user_schema_type" // TODO: naming?
)

func NewAddSchemaTypeUniqueConstraint(schemaType string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueSchemaType,
		schemaType,
		"Errors.UserSchema.Type.AlreadyExists") // TODO: i18n
}

func NewRemoveSchemaTypeUniqueConstraint(schemaType string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueSchemaType,
		schemaType,
	)
}

type CreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SchemaType             string                     `json:"type"`
	Schema                 map[string]any             `json:"schema"`
	PossibleAuthenticators []domain.AuthenticatorType `json:"possible_authenticators"`
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
	schema map[string]any,
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
