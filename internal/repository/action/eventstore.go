package action

import "github.com/zitadel/zitadel/v2/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedEventType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, ChangedEventType, ChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, DeactivatedEventType, DeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, ReactivatedEventType, ReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, RemovedEventType, RemovedEventMapper)
}
