package quota

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedEventType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, RemovedEventType, RemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationDueEventType, NotificationDueEventMapper).
		RegisterFilterEventMapper(AggregateType, NotifiedEventType, NotifiedEventMapper)
}
