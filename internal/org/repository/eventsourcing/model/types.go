package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	OrgAggregate       models.AggregateType = "org"
	OrgDomainAggregate models.AggregateType = "org.domain"
	OrgNameAggregate   models.AggregateType = "org.name"

	OrgAdded       models.EventType = "org.added"
	OrgChanged     models.EventType = "org.changed"
	OrgDeactivated models.EventType = "org.deactivated"
	OrgReactivated models.EventType = "org.reactivated"
	OrgRemoved     models.EventType = "org.removed"

	OrgNameReserved models.EventType = "org.name.reserved"
	OrgNameReleased models.EventType = "org.name.released"

	OrgDomainReserved models.EventType = "org.domain.reserved"
	OrgDomainReleased models.EventType = "org.domain.released"

	OrgMemberAdded   models.EventType = "org.member.added"
	OrgMemberChanged models.EventType = "org.member.changed"
	OrgMemberRemoved models.EventType = "org.member.removed"
)
