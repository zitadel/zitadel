package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	OrgAggregate models.AggregateType = "org"

	OrgAdded       models.EventType = "org.added"
	OrgChanged     models.EventType = "org.changed"
	OrgDeactivated models.EventType = "org.deactivated"
	OrgReactivated models.EventType = "org.reactivated"

	OrgMemberAdded   models.EventType = "org.member.added"
	OrgMemberChanged models.EventType = "org.member.changed"
	OrgMemberRemoved models.EventType = "org.member.removed"

	GrantMemberAdded   models.EventType = "org.grant.member.added"
	GrantMemberChanged models.EventType = "org.grant.member.changed"
	GrantMemberRemoved models.EventType = "org.grant.member.removed"
)
