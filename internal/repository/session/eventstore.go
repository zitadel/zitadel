package session

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	//es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
	es.RegisterFilterEventMapper(AggregateType, SetType, SetEventMapper).
		RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
