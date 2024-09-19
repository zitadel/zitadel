package target

import "github.com/zitadel/zitadel/v2/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, eventstore.GenericEventMapper[AddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ChangedEventType, eventstore.GenericEventMapper[ChangedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, eventstore.GenericEventMapper[RemovedEvent])
}
