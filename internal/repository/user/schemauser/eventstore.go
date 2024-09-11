package schemauser

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, CreatedType, eventstore.GenericEventMapper[CreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UpdatedType, eventstore.GenericEventMapper[UpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeletedType, eventstore.GenericEventMapper[DeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, LockedType, eventstore.GenericEventMapper[LockedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UnlockedType, eventstore.GenericEventMapper[UnlockedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ActivatedType, eventstore.GenericEventMapper[ActivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeactivatedType, eventstore.GenericEventMapper[DeactivatedEvent])
}
