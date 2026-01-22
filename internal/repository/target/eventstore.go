package target

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, eventstore.GenericEventMapper[AddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ChangedEventType, eventstore.GenericEventMapper[ChangedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, eventstore.GenericEventMapper[RemovedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, KeyAddedEventType, eventstore.GenericEventMapper[KeyAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, KeyActivatedEventType, eventstore.GenericEventMapper[KeyActivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, KeyDeactivatedEventType, eventstore.GenericEventMapper[KeyDeactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, KeyRemovedEventType, eventstore.GenericEventMapper[KeyRemovedEvent])
}
