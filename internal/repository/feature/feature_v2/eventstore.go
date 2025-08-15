package feature_v2

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SystemResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemLoginDefaultOrgEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemUserSchemaEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemTokenExchangeEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemImprovedPerformanceEventType, eventstore.GenericEventMapper[SetEvent[[]feature.ImprovedPerformanceType]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemOIDCSingleV1SessionTerminationEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemEnableBackChannelLogout, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemLoginVersion, eventstore.GenericEventMapper[SetEvent[*feature.LoginV2]])
	eventstore.RegisterFilterEventMapper(AggregateType, SystemPermissionCheckV2, eventstore.GenericEventMapper[SetEvent[bool]])

	eventstore.RegisterFilterEventMapper(AggregateType, InstanceResetEventType, eventstore.GenericEventMapper[ResetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLoginDefaultOrgEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceUserSchemaEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceTokenExchangeEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceImprovedPerformanceEventType, eventstore.GenericEventMapper[SetEvent[[]feature.ImprovedPerformanceType]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceDebugOIDCParentErrorEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceOIDCSingleV1SessionTerminationEventType, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceEnableBackChannelLogout, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceLoginVersion, eventstore.GenericEventMapper[SetEvent[*feature.LoginV2]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstancePermissionCheckV2, eventstore.GenericEventMapper[SetEvent[bool]])
	eventstore.RegisterFilterEventMapper(AggregateType, InstanceConsoleUseV2UserApi, eventstore.GenericEventMapper[SetEvent[bool]])
}
