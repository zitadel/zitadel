package debug_events

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, DebugAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ChangedEventType, DebugChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, DebugRemovedEventMapper)
}
