package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueOrgDomain                      = "org_domain"
	domainEventPrefix                    = orgEventTypePrefix + "domain."
	OrgDomainAddedEventType              = domainEventPrefix + "added"
	OrgDomainVerificationAddedEventType  = domainEventPrefix + "verification.added"
	OrgDomainVerificationFailedEventType = domainEventPrefix + "verification.failed"
	OrgDomainVerifiedEventType           = domainEventPrefix + "verified"
	OrgDomainPrimarySetEventType         = domainEventPrefix + "primary.set"
	OrgDomainRemovedEventType            = domainEventPrefix + "removed"
)

func NewAddOrgDomainUniqueConstraint(orgDomain string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueOrgDomain,
		orgDomain,
		"Errors.Org.Domain.AlreadyExists")
}

func NewRemoveOrgDomainUniqueConstraint(orgDomain string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueOrgDomain,
		orgDomain)
}

type DomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainAddedEvent) Payload() interface{} {
	return e
}

func (e *DomainAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainAddedEvent {
	return &DomainAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainAddedEventType,
		),
		Domain: domain,
	}
}

func DomainAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainAdded := &DomainAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-GBr52", "unable to unmarshal org domain added")
	}

	return orgDomainAdded, nil
}

type DomainVerificationAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain         string                         `json:"domain,omitempty"`
	ValidationType domain.OrgDomainValidationType `json:"validationType,omitempty"`
	ValidationCode *crypto.CryptoValue            `json:"validationCode,omitempty"`
}

func (e *DomainVerificationAddedEvent) Payload() interface{} {
	return e
}

func (e *DomainVerificationAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainVerificationAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	domain string,
	validationType domain.OrgDomainValidationType,
	validationCode *crypto.CryptoValue) *DomainVerificationAddedEvent {
	return &DomainVerificationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainVerificationAddedEventType,
		),
		Domain:         domain,
		ValidationType: validationType,
		ValidationCode: validationCode,
	}
}

func DomainVerificationAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainVerificationAdded := &DomainVerificationAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainVerificationAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-NRN32", "unable to unmarshal org domain verification added")
	}

	return orgDomainVerificationAdded, nil
}

type DomainVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainVerificationFailedEvent) Payload() interface{} {
	return e
}

func (e *DomainVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainVerificationFailedEvent {
	return &DomainVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainVerificationFailedEventType,
		),
		Domain: domain,
	}
}

func DomainVerificationFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainVerificationFailed := &DomainVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainVerificationFailed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-Bhm37", "unable to unmarshal org domain verification failed")
	}

	return orgDomainVerificationFailed, nil
}

type DomainVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainVerifiedEvent) Payload() interface{} {
	return e
}

func (e *DomainVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddOrgDomainUniqueConstraint(e.Domain)}
}

func NewDomainVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainVerifiedEvent {
	return &DomainVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainVerifiedEventType,
		),
		Domain: domain,
	}
}

func DomainVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainVerified := &DomainVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainVerified)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-BFSwt", "unable to unmarshal org domain verified")
	}

	return orgDomainVerified, nil
}

type DomainPrimarySetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainPrimarySetEvent) Payload() interface{} {
	return e
}

func (e *DomainPrimarySetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainPrimarySetEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainPrimarySetEvent {
	return &DomainPrimarySetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainPrimarySetEventType,
		),
		Domain: domain,
	}
}

func DomainPrimarySetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainPrimarySet := &DomainPrimarySetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainPrimarySet)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-N5787", "unable to unmarshal org domain primary set")
	}

	return orgDomainPrimarySet, nil
}

type DomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain     string `json:"domain,omitempty"`
	isVerified bool
}

func (e *DomainRemovedEvent) Payload() interface{} {
	return e
}

func (e *DomainRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if !e.isVerified {
		return nil
	}
	return []*eventstore.UniqueConstraint{NewRemoveOrgDomainUniqueConstraint(e.Domain)}
}

func NewDomainRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string, verified bool) *DomainRemovedEvent {
	return &DomainRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainRemovedEventType,
		),
		Domain:     domain,
		isVerified: verified,
	}
}

func DomainRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgDomainRemoved := &DomainRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgDomainRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-BngB2", "unable to unmarshal org domain removed")
	}

	return orgDomainRemoved, nil
}
