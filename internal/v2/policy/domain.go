package policy

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type domainPolicyAddedPayload struct {
	UserLoginMustBeDomain                  bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

type DomainPolicyAddedEvent domainPolicyAddedEvent
type domainPolicyAddedEvent = eventstore.Event[domainPolicyAddedPayload]

func DomainPolicyAddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyAddedEvent, error) {
	event, err := eventstore.EventFromStorage[domainPolicyAddedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*DomainPolicyAddedEvent)(event), nil
}

func (e *DomainPolicyAddedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "policy.domain.added")
}

type domainPolicyChangedPayload struct {
	UserLoginMustBeDomain                  *bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     *bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain *bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

type DomainPolicyChangedEvent domainPolicyChangedEvent
type domainPolicyChangedEvent = eventstore.Event[domainPolicyChangedPayload]

func DomainPolicyChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyChangedEvent, error) {
	event, err := eventstore.EventFromStorage[domainPolicyChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*DomainPolicyChangedEvent)(event), nil
}

func (e *DomainPolicyChangedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "policy.domain.changed")
}

type DomainPolicyRemovedEvent domainPolicyRemovedEvent
type domainPolicyRemovedEvent = eventstore.Event[struct{}]

func DomainPolicyRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[domainPolicyRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*DomainPolicyRemovedEvent)(event), nil
}

func (e *DomainPolicyRemovedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "policy.domain.removed")
}
