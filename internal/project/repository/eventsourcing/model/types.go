package model

import "github.com/caos/zitadel/internal/eventstore/v1/models"

const (
	ProjectAggregate models.AggregateType = "project"

	ProjectAdded       models.EventType = "project.added"
	ProjectChanged     models.EventType = "project.changed"
	ProjectDeactivated models.EventType = "project.deactivated"
	ProjectReactivated models.EventType = "project.reactivated"
	ProjectRemoved     models.EventType = "project.removed"

	ProjectMemberAdded          models.EventType = "project.member.added"
	ProjectMemberChanged        models.EventType = "project.member.changed"
	ProjectMemberRemoved        models.EventType = "project.member.removed"
	ProjectMemberCascadeRemoved models.EventType = "project.member.cascade.removed"

	ProjectRoleAdded   models.EventType = "project.role.added"
	ProjectRoleChanged models.EventType = "project.role.changed"
	ProjectRoleRemoved models.EventType = "project.role.removed"

	ProjectGrantAdded          models.EventType = "project.grant.added"
	ProjectGrantChanged        models.EventType = "project.grant.changed"
	ProjectGrantRemoved        models.EventType = "project.grant.removed"
	ProjectGrantDeactivated    models.EventType = "project.grant.deactivated"
	ProjectGrantReactivated    models.EventType = "project.grant.reactivated"
	ProjectGrantCascadeChanged models.EventType = "project.grant.cascade.changed"

	ProjectGrantMemberAdded          models.EventType = "project.grant.member.added"
	ProjectGrantMemberChanged        models.EventType = "project.grant.member.changed"
	ProjectGrantMemberRemoved        models.EventType = "project.grant.member.removed"
	ProjectGrantMemberCascadeRemoved models.EventType = "project.grant.member.cascade.removed"

	ApplicationAdded       models.EventType = "project.application.added"
	ApplicationChanged     models.EventType = "project.application.changed"
	ApplicationRemoved     models.EventType = "project.application.removed"
	ApplicationDeactivated models.EventType = "project.application.deactivated"
	ApplicationReactivated models.EventType = "project.application.reactivated"

	OIDCConfigAdded                models.EventType = "project.application.config.oidc.added"
	OIDCConfigChanged              models.EventType = "project.application.config.oidc.changed"
	OIDCConfigSecretChanged        models.EventType = "project.application.config.oidc.secret.changed"
	OIDCClientSecretCheckSucceeded models.EventType = "project.application.oidc.secret.check.succeeded"
	OIDCClientSecretCheckFailed    models.EventType = "project.application.oidc.secret.check.failed"

	APIConfigAdded         models.EventType = "project.application.config.api.added"
	APIConfigChanged       models.EventType = "project.application.config.api.changed"
	APIConfigSecretChanged models.EventType = "project.application.config.api.secret.changed"

	ClientKeyAdded   models.EventType = "project.application.oidc.key.added"
	ClientKeyRemoved models.EventType = "project.application.oidc.key.removed"
)
