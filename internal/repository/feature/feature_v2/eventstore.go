package feature_v2

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SystemResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemDefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLoginDefaultOrgEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
}
