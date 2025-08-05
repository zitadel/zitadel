package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, eventstore.GenericEventMapper[AddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ApprovedEventType, eventstore.GenericEventMapper[ApprovedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, CanceledEventType, eventstore.GenericEventMapper[CanceledEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DoneEventType, eventstore.GenericEventMapper[DoneEvent])
}
