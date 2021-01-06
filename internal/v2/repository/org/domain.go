package org

const (
	domainEventPrefix           = orgEventTypePrefix + "domain."
	OrgDomainAdded              = domainEventPrefix + "added"
	OrgDomainVerificationAdded  = domainEventPrefix + "verification.added"
	OrgDomainVerificationFailed = domainEventPrefix + "verification.failed"
	OrgDomainVerified           = domainEventPrefix + "verified"
	OrgDomainRemoved            = domainEventPrefix + "removed"
	OrgDomainPrimarySet         = domainEventPrefix + "primary.set"
)
