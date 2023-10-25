package limits

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper).
		RegisterFilterEventMapper(AggregateType, ResetEventType, ResetEventMapper)
}
