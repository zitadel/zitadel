package feature

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, DefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[Boolean]])
}
