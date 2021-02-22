package project

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(ProjectAddedType, ProjectAddedEventMapper).
		RegisterFilterEventMapper(ProjectChangedType, ProjectChangeEventMapper).
		RegisterFilterEventMapper(ProjectDeactivatedType, ProjectDeactivatedEventMapper).
		RegisterFilterEventMapper(ProjectReactivatedType, ProjectReactivatedEventMapper).
		RegisterFilterEventMapper(ProjectRemovedType, ProjectRemovedEventMapper).
		RegisterFilterEventMapper(MemberAddedType, MemberAddedEventMapper).
		RegisterFilterEventMapper(MemberChangedType, MemberChangedEventMapper).
		RegisterFilterEventMapper(MemberRemovedType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(RoleAddedType, RoleAddedEventMapper).
		RegisterFilterEventMapper(RoleChangedType, RoleChangedEventMapper).
		RegisterFilterEventMapper(RoleRemovedType, RoleRemovedEventMapper).
		RegisterFilterEventMapper(GrantAddedType, GrantAddedEventMapper).
		RegisterFilterEventMapper(GrantChangedType, GrantChangedEventMapper).
		RegisterFilterEventMapper(GrantCascadeChangedType, GrantCascadeChangedEventMapper).
		RegisterFilterEventMapper(GrantDeactivatedType, GrantDeactivateEventMapper).
		RegisterFilterEventMapper(GrantReactivatedType, GrantReactivatedEventMapper).
		RegisterFilterEventMapper(GrantRemovedType, GrantRemovedEventMapper).
		RegisterFilterEventMapper(GrantMemberAddedType, GrantMemberAddedEventMapper).
		RegisterFilterEventMapper(GrantMemberChangedType, GrantMemberChangedEventMapper).
		RegisterFilterEventMapper(GrantMemberRemovedType, GrantMemberRemovedEventMapper).
		RegisterFilterEventMapper(ApplicationAddedType, ApplicationAddedEventMapper).
		RegisterFilterEventMapper(ApplicationChangedType, ApplicationChangedEventMapper).
		RegisterFilterEventMapper(ApplicationRemovedType, ApplicationRemovedEventMapper).
		RegisterFilterEventMapper(ApplicationDeactivatedType, ApplicationDeactivatedEventMapper).
		RegisterFilterEventMapper(ApplicationReactivatedType, ApplicationReactivatedEventMapper).
		RegisterFilterEventMapper(OIDCConfigAddedType, OIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(OIDCConfigChangedType, OIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(OIDCConfigSecretChangedType, OIDCConfigSecretChangedEventMapper).
		RegisterFilterEventMapper(OIDCClientSecretCheckSucceededType, OIDCConfigSecretCheckSucceededEventMapper).
		RegisterFilterEventMapper(OIDCClientSecretCheckFailedType, OIDCConfigSecretCheckFailedEventMapper).
		RegisterFilterEventMapper(APIConfigAddedType, APIConfigAddedEventMapper).
		RegisterFilterEventMapper(APIConfigChangedType, APIConfigChangedEventMapper).
		RegisterFilterEventMapper(APIConfigSecretChangedType, APIConfigSecretChangedEventMapper).
		RegisterFilterEventMapper(ApplicationKeyAddedEventType, ApplicationKeyAddedEventMapper).
		RegisterFilterEventMapper(ApplicationKeyRemovedEventType, ApplicationKeyRemovedEventMapper)
}
