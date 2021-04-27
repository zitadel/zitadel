package org

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func NewAddOrgDomainUniqueConstraint(orgDomain string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueOrgDomain,
		orgDomain,
		"Errors.Org.Domain.AlreadyExists")
}

func NewRemoveOrgDomainUniqueConstraint(orgDomain string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueOrgDomain,
		orgDomain)
}

type DomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainAddedEvent) Data() interface{} {
	return e
}

func (e *DomainAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *DomainAddedEvent) Assets() []*eventstore.Asset {
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

func DomainAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainAdded := &DomainAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-GBr52", "unable to unmarshal org domain added")
	}

	return orgDomainAdded, nil
}

type DomainVerificationAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain         string                         `json:"domain,omitempty"`
	ValidationType domain.OrgDomainValidationType `json:"validationType,omitempty"`
	ValidationCode *crypto.CryptoValue            `json:"validationCode,omitempty"`
}

func (e *DomainVerificationAddedEvent) Data() interface{} {
	return e
}

func (e *DomainVerificationAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *DomainVerificationAddedEvent) Assets() []*eventstore.Asset {
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

func DomainVerificationAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainVerificationAdded := &DomainVerificationAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainVerificationAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-NRN32", "unable to unmarshal org domain verification added")
	}

	return orgDomainVerificationAdded, nil
}

type DomainVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainVerificationFailedEvent) Data() interface{} {
	return e
}

func (e *DomainVerificationFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *DomainVerificationFailedEvent) Assets() []*eventstore.Asset {
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

func DomainVerificationFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainVerificationFailed := &DomainVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainVerificationFailed)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-Bhm37", "unable to unmarshal org domain verification failed")
	}

	return orgDomainVerificationFailed, nil
}

type DomainVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainVerifiedEvent) Data() interface{} {
	return e
}

func (e *DomainVerifiedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddOrgDomainUniqueConstraint(e.Domain)}
}

func (e *DomainVerifiedEvent) Assets() []*eventstore.Asset {
	return nil
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

func DomainVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainVerified := &DomainVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-BFSwt", "unable to unmarshal org domain verified")
	}

	return orgDomainVerified, nil
}

type DomainPrimarySetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainPrimarySetEvent) Data() interface{} {
	return e
}

func (e *DomainPrimarySetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *DomainPrimarySetEvent) Assets() []*eventstore.Asset {
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

func DomainPrimarySetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainPrimarySet := &DomainPrimarySetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainPrimarySet)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-N5787", "unable to unmarshal org domain primary set")
	}

	return orgDomainPrimarySet, nil
}

type DomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain     string `json:"domain,omitempty"`
	isVerified bool
}

func (e *DomainRemovedEvent) Data() interface{} {
	return e
}

func (e *DomainRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	if !e.isVerified {
		return nil
	}
	return []*eventstore.EventUniqueConstraint{NewRemoveOrgDomainUniqueConstraint(e.Domain)}
}

func (e *DomainRemovedEvent) Assets() []*eventstore.Asset {
	return nil
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

func DomainRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainRemoved := &DomainRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-BngB2", "unable to unmarshal org domain removed")
	}

	return orgDomainRemoved, nil
}
