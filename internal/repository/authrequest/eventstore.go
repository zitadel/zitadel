package authrequest

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper)
	//RegisterFilterEventMapper(AggregateType, UserCheckedType, UserCheckedEventMapper).
	//RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
