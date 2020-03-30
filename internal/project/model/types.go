package model

const (
	AddedProject       = "project.added"
	ChangedProject     = "project.changed"
	DeactivatedProject = "project.deactivated"
	ReactivatedProject = "project.reactivated"

	AddedMember   = "project.member.added"
	ChangedMember = "project.member.changed"
	RemovedMember = "project.member.removed"

	AddedRole   = "project.role.added"
	RemovedRole = "project.role.removed"

	AddedProjectGrant       = "project.grant.added"
	ChangedProjectGrant     = "project.grant.changed"
	DeactivatedProjectGrant = "project.grant.deactivated"
	ReactivatedProjectGrant = "project.grant.reactivated"

	AddedGrantMember   = "project.grant.member.added"
	ChangedGrantMember = "project.grant.member.changed"
	RemovedGrantMember = "project.grant.member.removed"

	AddedApplication       = "project.application.added"
	ChangedApplication     = "project.application.changed"
	DeactivatedApplication = "project.application.deactivated"
	ReactivatedApplication = "project.application.reactivated"

	AddedOIDCConfig         = "project.application.config.oidc.added"
	ChangedOIDCConfig       = "project.application.config.oidc.changed"
	ChangedOIDCConfigSecret = "project.application.config.oidc.secret.changed"
)

const (
	ProjectAggregate = "project"
)
