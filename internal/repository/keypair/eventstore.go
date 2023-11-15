package keypair

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedEventType, AddedEventMapper)
	es.RegisterFilterEventMapper(AggregateType, AddedCertificateEventType, AddedCertificateEventMapper)
}
