package project

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, ProjectAddedType, ProjectAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, ProjectChangedType, ProjectChangeEventMapper).
		RegisterFilterEventMapper(AggregateType, ProjectDeactivatedType, ProjectDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, ProjectReactivatedType, ProjectReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, ProjectRemovedType, ProjectRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberAddedType, MemberAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberChangedType, MemberChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberRemovedType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberCascadeRemovedType, MemberCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, RoleAddedType, RoleAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, RoleChangedType, RoleChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, RoleRemovedType, RoleRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantAddedType, GrantAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantChangedType, GrantChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantCascadeChangedType, GrantCascadeChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantDeactivatedType, GrantDeactivateEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantReactivatedType, GrantReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantRemovedType, GrantRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantMemberAddedType, GrantMemberAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantMemberChangedType, GrantMemberChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantMemberRemovedType, GrantMemberRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, GrantMemberCascadeRemovedType, GrantMemberCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationAddedType, ApplicationAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationChangedType, ApplicationChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationRemovedType, ApplicationRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationDeactivatedType, ApplicationDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationReactivatedType, ApplicationReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCConfigAddedType, OIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCConfigChangedType, OIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCConfigSecretChangedType, OIDCConfigSecretChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCClientSecretCheckSucceededType, OIDCConfigSecretCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCClientSecretCheckFailedType, OIDCConfigSecretCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, APIConfigAddedType, APIConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, APIConfigChangedType, APIConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, APIConfigSecretChangedType, APIConfigSecretChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationKeyAddedEventType, ApplicationKeyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, ApplicationKeyRemovedEventType, ApplicationKeyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, SAMLConfigAddedType, SAMLConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SAMLConfigChangedType, SAMLConfigChangedEventMapper)
}
