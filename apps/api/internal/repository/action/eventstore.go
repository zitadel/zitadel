package action

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, AddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ChangedEventType, ChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, DeactivatedEventType, DeactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ReactivatedEventType, ReactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, RemovedEventMapper)
}
