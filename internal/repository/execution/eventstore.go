package target

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SetRequestEventType, SetRequestEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SetResponseEventType, SetResponseEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SetFunctionEventType, SetFunctionEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventEventType, SetEventEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedRequestEventType, RemovedRequestEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedResponseEventType, RemovedResponseEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedFunctionEventType, RemovedFunctionEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RemovedEventEventType, RemovedEventEventMapper)
}
