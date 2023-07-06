package authrequest

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SessionLinkedType, SessionLinkedEventMapper).
		RegisterFilterEventMapper(AggregateType, CodeAddedType, CodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, CodeExchangedType, CodeExchangedEventMapper).
		RegisterFilterEventMapper(AggregateType, FailedType, FailedEventMapper).
		RegisterFilterEventMapper(AggregateType, SucceededType, SucceededEventMapper)
}
