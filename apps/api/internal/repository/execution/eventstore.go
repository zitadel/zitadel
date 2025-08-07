package execution

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventType, eventstore.GenericEventMapper[SetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventV2Type, eventstore.GenericEventMapper[SetEventV2])
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventType, eventstore.GenericEventMapper[RemovedEvent])
}
