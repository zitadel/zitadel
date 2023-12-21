package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedEventType, eventstore.GenericEventMapper[AddedEvent]).
		RegisterFilterEventMapper(AggregateType, ApprovedEventType, eventstore.GenericEventMapper[ApprovedEvent]).
		RegisterFilterEventMapper(AggregateType, CanceledEventType, eventstore.GenericEventMapper[CanceledEvent])
}
