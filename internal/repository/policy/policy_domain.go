package policy

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DomainPolicyAddedEventType   = "policy.domain.added"
	DomainPolicyChangedEventType = "policy.domain.changed"
	DomainPolicyRemovedEventType = "policy.domain.removed"
)

type DomainPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain                  bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

func (e *DomainPolicyAddedEvent) Payload() interface{} {
	return e
}

func (e *DomainPolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainPolicyAddedEvent(
	base *eventstore.BaseEvent,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) *DomainPolicyAddedEvent {

	return &DomainPolicyAddedEvent{
		BaseEvent:                              *base,
		UserLoginMustBeDomain:                  userLoginMustBeDomain,
		ValidateOrgDomains:                     validateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: smtpSenderAddressMatchesInstanceDomain,
	}
}

func DomainPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DomainPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-TvSmA", "unable to unmarshal policy")
	}

	return e, nil
}

type DomainPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain                  *bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     *bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain *bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

func (e *DomainPolicyChangedEvent) Payload() interface{} {
	return e
}

func (e *DomainPolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []DomainPolicyChanges,
) (*DomainPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-DAf3h", "Errors.NoChangesFound")
	}
	changeEvent := &DomainPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type DomainPolicyChanges func(*DomainPolicyChangedEvent)

func ChangeUserLoginMustBeDomain(userLoginMustBeDomain bool) func(*DomainPolicyChangedEvent) {
	return func(e *DomainPolicyChangedEvent) {
		e.UserLoginMustBeDomain = &userLoginMustBeDomain
	}
}

func ChangeValidateOrgDomains(validateOrgDomain bool) func(*DomainPolicyChangedEvent) {
	return func(e *DomainPolicyChangedEvent) {
		e.ValidateOrgDomains = &validateOrgDomain
	}
}

func ChangeSMTPSenderAddressMatchesInstanceDomain(smtpSenderAddressMatchesInstanceDomain bool) func(*DomainPolicyChangedEvent) {
	return func(e *DomainPolicyChangedEvent) {
		e.SMTPSenderAddressMatchesInstanceDomain = &smtpSenderAddressMatchesInstanceDomain
	}
}

func DomainPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DomainPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-0Pl9d", "unable to unmarshal policy")
	}

	return e, nil
}

type DomainPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DomainPolicyRemovedEvent) Payload() interface{} {
	return nil
}

func (e *DomainPolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainPolicyRemovedEvent(base *eventstore.BaseEvent) *DomainPolicyRemovedEvent {
	return &DomainPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func DomainPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &DomainPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
