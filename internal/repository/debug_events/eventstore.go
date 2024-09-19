package debug_events

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, DebugAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ChangedEventType, DebugChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, DebugRemovedEventMapper)
}
