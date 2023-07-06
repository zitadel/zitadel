package milestone

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, PushedEventType, PushedEventMapper)
}
