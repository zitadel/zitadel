package authenticator

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameCreatedType, eventstore.GenericEventMapper[UsernameCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UsernameDeletedType, eventstore.GenericEventMapper[UsernameDeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCreatedType, eventstore.GenericEventMapper[PasswordCreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordDeletedType, eventstore.GenericEventMapper[PasswordDeletedEvent])
}
