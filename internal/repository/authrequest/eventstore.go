package authrequest

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, CodeAddedType, CodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SessionLinkedType, SessionLinkedEventMapper).
		RegisterFilterEventMapper(AggregateType, FailedType, FailedEventMapper)
}
