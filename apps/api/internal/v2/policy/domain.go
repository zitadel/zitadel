package policy

import "github.com/zitadel/zitadel/internal/v2/eventstore"

const DomainPolicyAddedTypeSuffix = "policy.domain.added"

type DomainPolicyAddedPayload struct {
	UserLoginMustBeDomain                  bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

const DomainPolicyChangedTypeSuffix = "policy.domain.changed"

type DomainPolicyChangedPayload struct {
	UserLoginMustBeDomain                  *bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     *bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain *bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

const DomainPolicyRemovedTypeSuffix = "policy.domain.removed"

type DomainPolicyRemovedPayload eventstore.EmptyPayload
