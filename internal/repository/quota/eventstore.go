package quota

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	// AddedEventType is not emitted anymore.
	// For ease of use, old events are directly mapped to SetEvent.
	// This works, because the data structures are compatible.
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, SetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, RemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, NotificationDueEventType, NotificationDueEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, NotifiedEventType, NotifiedEventMapper)
}
