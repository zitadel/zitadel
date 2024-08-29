package schemauser

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, CreatedType, eventstore.GenericEventMapper[CreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UpdatedType, eventstore.GenericEventMapper[UpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeletedType, eventstore.GenericEventMapper[DeletedEvent])
}
