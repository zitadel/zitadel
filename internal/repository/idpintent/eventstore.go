package idpintent

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, StartedEventType, StartedEventMapper).
		RegisterFilterEventMapper(AggregateType, SucceededEventType, SucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, SAMLSucceededEventType, SAMLSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, SAMLRequestEventType, SAMLRequestEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPSucceededEventType, LDAPSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, FailedEventType, FailedEventMapper)
}
