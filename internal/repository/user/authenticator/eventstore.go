package authenticator

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameCreatedType, eventstore.GenericEventMapper[UsernameCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameDeletedType, eventstore.GenericEventMapper[UsernameDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCreatedType, eventstore.GenericEventMapper[PasswordCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCodeAddedType, eventstore.GenericEventMapper[PasswordCodeAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCodeSentType, eventstore.GenericEventMapper[PasswordCodeSentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordDeletedType, eventstore.GenericEventMapper[PasswordDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PublicKeyCreatedType, eventstore.GenericEventMapper[PublicKeyCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PublicKeyDeletedType, eventstore.GenericEventMapper[PublicKeyDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PATCreatedType, eventstore.GenericEventMapper[PATCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PATDeletedType, eventstore.GenericEventMapper[PATDeletedEvent])
}
