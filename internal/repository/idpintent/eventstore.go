package idpintent

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, StartedEventType, StartedEventMapper).
		RegisterFilterEventMapper(AggregateType, SucceededEventType, SucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPSucceededEventType, LDAPSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, FailedEventType, FailedEventMapper)
}
