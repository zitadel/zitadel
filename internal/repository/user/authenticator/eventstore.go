package authenticator

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameCreatedType, eventstore.GenericEventMapper[UsernameCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameDeletedType, eventstore.GenericEventMapper[UsernameDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCreatedType, eventstore.GenericEventMapper[PasswordCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCodeAddedType, eventstore.GenericEventMapper[PasswordCodeAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordDeletedType, eventstore.GenericEventMapper[PasswordDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, JWTCreatedType, eventstore.GenericEventMapper[JWTCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, JWTDeletedType, eventstore.GenericEventMapper[JWTDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PATCreatedType, eventstore.GenericEventMapper[PATCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PATDeletedType, eventstore.GenericEventMapper[PATDeletedEvent])
}
