package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	OrgAggregate       models.AggregateType = "org"
	OrgDomainAggregate models.AggregateType = "org.domain"
	OrgNameAggregate   models.AggregateType = "org.name"

	OrgAdded            models.EventType = "org.added"
	OrgChanged          models.EventType = "org.changed"
	OrgDeactivated      models.EventType = "org.deactivated"
	OrgReactivated      models.EventType = "org.reactivated"
	OrgRemoved          models.EventType = "org.removed"
	OrgDomainAdded      models.EventType = "org.domain.added"
	OrgDomainVerified   models.EventType = "org.domain.verified"
	OrgDomainRemoved    models.EventType = "org.domain.removed"
	OrgDomainPrimarySet models.EventType = "org.domain.primary.set"

	OrgNameReserved models.EventType = "org.name.reserved"
	OrgNameReleased models.EventType = "org.name.released"

	OrgDomainReserved models.EventType = "org.domain.reserved"
	OrgDomainReleased models.EventType = "org.domain.released"

	OrgMemberAdded   models.EventType = "org.member.added"
	OrgMemberChanged models.EventType = "org.member.changed"
	OrgMemberRemoved models.EventType = "org.member.removed"

	OrgIamPolicyAdded   models.EventType = "org.iam.policy.added"
	OrgIamPolicyChanged models.EventType = "org.iam.policy.changed"
	OrgIamPolicyRemoved models.EventType = "org.iam.policy.removed"
)
