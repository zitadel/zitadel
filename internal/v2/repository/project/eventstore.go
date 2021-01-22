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
		RegisterFilterEventMapper(ProjectMemberAddedType, ProjectMemberAddedEventMapper).
		RegisterFilterEventMapper(ProjectMemberChangedType, ProjectMemberChangedEventMapper).
		RegisterFilterEventMapper(ProjectMemberRemovedType, ProjectMemberRemovedEventMapper).
		RegisterFilterEventMapper(RoleAddedType, RoleAddedEventMapper).
		RegisterFilterEventMapper(RoleChangedType, RoleChangedEventMapper).
		RegisterFilterEventMapper(RoleRemovedType, RoleRemovedEventMapper).
		RegisterFilterEventMapper(GrantAddedType, GrantAddedEventMapper).
		RegisterFilterEventMapper(GrantChangedType, GrantChangedEventMapper).
		RegisterFilterEventMapper(GrantCascadeChangedType, GrantChangedEventMapper).
		RegisterFilterEventMapper(GrantDeactivatedType, GrantDeactivateEventMapper).
		RegisterFilterEventMapper(GrantReactivatedType, GrantReactivatedEventMapper).
		RegisterFilterEventMapper(GrantRemovedType, GrantRemovedEventMapper).
		RegisterFilterEventMapper(ProjectGrantMemberAddedType, ProjectGrantMemberAddedEventMapper).
		RegisterFilterEventMapper(ProjectGrantMemberChangedType, ProjectGrantMemberChangedEventMapper).
		RegisterFilterEventMapper(ProjectGrantMemberRemovedType, ProjectGrantMemberChangedEventMapper).
		RegisterFilterEventMapper(ApplicationAddedType, ApplicationAddedEventMapper).
		RegisterFilterEventMapper(ApplicationChangedType, ApplicationAddedEventMapper).
		RegisterFilterEventMapper(ApplicationRemovedType, ApplicationRemovedEventMapper).
		RegisterFilterEventMapper(ApplicationDeactivatedType, ApplicationDeactivatedEventMapper).
		RegisterFilterEventMapper(ApplicationReactivatedType, ApplicationReactivatedEventMapper).
		RegisterFilterEventMapper(OIDCConfigAddedType, OIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(OIDCConfigChangedType, OIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(OIDCConfigSecretChangedType, OIDCConfigSecretChangedEventMapper).
		RegisterFilterEventMapper(OIDCClientSecretCheckSucceededType, OIDCConfigSecretCheckSucceededEventMapper).
		RegisterFilterEventMapper(OIDCClientSecretCheckFailedType, OIDCConfigSecretCheckFailedEventMapper)
}
