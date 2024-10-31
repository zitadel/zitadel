package feature_v2

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SystemResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemLoginDefaultOrgEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemUserSchemaEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemTokenExchangeEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemActionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemImprovedPerformanceEventType, eventstore.GenericEventMapper[SetEvent[[]feature.ImprovedPerformanceType]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemOIDCSingleV1SessionTerminationEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemDisableUserTokenEvent, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemEnableBackChannelLogout, eventstore.GenericEventMapper[SetEvent[bool]])

	eventstore.RegisterFilterEventMapper(AggregateType, InstanceResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLoginDefaultOrgEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceTriggerIntrospectionProjectionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLegacyIntrospectionEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceUserSchemaEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceTokenExchangeEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceActionsEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceImprovedPerformanceEventType, eventstore.GenericEventMapper[SetEvent[[]feature.ImprovedPerformanceType]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceWebKeyEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceDebugOIDCParentErrorEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceOIDCSingleV1SessionTerminationEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceDisableUserTokenEvent, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceEnableBackChannelLogout, eventstore.GenericEventMapper[SetEvent[bool]])
}
