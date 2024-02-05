package feature_v2

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, SystemResetEventType, eventstore.GenericEventMapper[ResetEvent])
	es.RegisterFilterEventMapper(AggregateType, SystemDefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	es.RegisterFilterEventMapper(AggregateType, SystemTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	es.RegisterFilterEventMapper(AggregateType, SystemLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	es.RegisterFilterEventMapper(AggregateType, InstanceResetEventType, eventstore.GenericEventMapper[ResetEvent])
	es.RegisterFilterEventMapper(AggregateType, InstanceDefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	es.RegisterFilterEventMapper(AggregateType, InstanceTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	es.RegisterFilterEventMapper(AggregateType, InstanceLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
}
