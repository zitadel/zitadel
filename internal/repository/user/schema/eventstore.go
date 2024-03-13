package schema

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, CreatedType, eventstore.GenericEventMapper[CreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UpdatedType, eventstore.GenericEventMapper[UpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeactivatedType, eventstore.GenericEventMapper[DeactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ReactivatedType, eventstore.GenericEventMapper[ReactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeletedType, eventstore.GenericEventMapper[DeletedEvent])
}
