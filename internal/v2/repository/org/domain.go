package org

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	domainEventPrefix           = orgEventTypePrefix + "domain."
	OrgDomainAdded              = domainEventPrefix + "added"
	OrgDomainVerificationAdded  = domainEventPrefix + "verification.added"
	OrgDomainVerificationFailed = domainEventPrefix + "verification.failed"
	OrgDomainVerified           = domainEventPrefix + "verified"
	OrgDomainPrimarySet         = domainEventPrefix + "primary.set"
	OrgDomainRemoved            = domainEventPrefix + "removed"
)

type DomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainAddedEvent) Data() interface{} {
	return e
}

func NewDomainAddedEvent(ctx context.Context, domain string) *DomainAddedEvent {
	return &DomainAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainAdded,
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
		return nil, errors.ThrowInternal(err, "USER-GBr52", "unable to unmarshal org domain added")
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

func NewDomainVerificationAddedEvent(
	ctx context.Context,
	domain string,
	validationType domain.OrgDomainValidationType,
	validationCode *crypto.CryptoValue) *DomainVerificationAddedEvent {
	return &DomainVerificationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainVerificationAdded,
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
		return nil, errors.ThrowInternal(err, "USER-NRN32", "unable to unmarshal org domain verification added")
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

func NewDomainVerificationFailedEvent(ctx context.Context, domain string) *DomainVerificationFailedEvent {
	return &DomainVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainVerificationFailed,
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
		return nil, errors.ThrowInternal(err, "USER-Bhm37", "unable to unmarshal org domain verification failed")
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

func NewDomainVerifiedEvent(ctx context.Context, domain string) *DomainVerifiedEvent {
	return &DomainVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainVerified,
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
		return nil, errors.ThrowInternal(err, "USER-BFSwt", "unable to unmarshal org domain verified")
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

func NewDomainPrimarySetEvent(ctx context.Context, domain string) *DomainPrimarySetEvent {
	return &DomainPrimarySetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainPrimarySet,
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
		return nil, errors.ThrowInternal(err, "USER-N5787", "unable to unmarshal org domain primary set")
	}

	return orgDomainPrimarySet, nil
}

type DomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainRemovedEvent) Data() interface{} {
	return e
}

func NewDomainRemovedEvent(ctx context.Context, domain string) *DomainRemovedEvent {
	return &DomainRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDomainRemoved,
		),
		Domain: domain,
	}
}

func DomainRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgDomainRemoved := &DomainRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-BngB2", "unable to unmarshal org domain removed")
	}

	return orgDomainRemoved, nil
}
