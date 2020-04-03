package model

const (
	ProjectAggregate = "project"

	ProjectAdded       = "project.added"
	ProjectChanged     = "project.changed"
	ProjectDeactivated = "project.deactivated"
	ProjectReactivated = "project.reactivated"

	ProjectMemberAdded   = "project.member.added"
	ProjectMemberChanged = "project.member.changed"
	ProjectMemberRemoved = "project.member.removed"

	ProjectRoleAdded   = "project.role.added"
	ProjectRoleRemoved = "project.role.removed"

	ProjectGrantAdded       = "project.grant.added"
	ProjectGrantChanged     = "project.grant.changed"
	ProjectGrantDeactivated = "project.grant.deactivated"
	ProjectGrantReactivated = "project.grant.reactivated"

	GrantMemberAdded   = "project.grant.member.added"
	GrantMemberChanged = "project.grant.member.changed"
	GrantMemberRemoved = "project.grant.member.removed"

	ApplicationAdded       = "project.application.added"
	ApplicationChanged     = "project.application.changed"
	ApplicationDeactivated = "project.application.deactivated"
	ApplicationReactivated = "project.application.reactivated"

	OIDCConfigAdded         = "project.application.config.oidc.added"
	OIDCConfigChanged       = "project.application.config.oidc.changed"
	OIDCConfigSecretChanged = "project.application.config.oidc.secret.changed"
)
