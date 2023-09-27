package quota

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	// AddedEventType is not emitted anymore.
	// For ease of use, old events are directly mapped to SetEvent.
	// This works, because the data structures are compatible.
	es.RegisterFilterEventMapper(AggregateType, AddedEventType, SetEventMapper).
		RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper).
		RegisterFilterEventMapper(AggregateType, RemovedEventType, RemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationDueEventType, NotificationDueEventMapper).
		RegisterFilterEventMapper(AggregateType, NotifiedEventType, NotifiedEventMapper)
}
