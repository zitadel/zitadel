package limits

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper).
		RegisterFilterEventMapper(AggregateType, ResetEventType, ResetEventMapper)
}
