package instance

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	UniqueInstanceDomain              = "instance_domain"
	domainEventPrefix                 = instanceEventTypePrefix + "domain."
	InstanceDomainAddedEventType      = domainEventPrefix + "added"
	InstanceDomainPrimarySetEventType = domainEventPrefix + "primary.set"
	InstanceDomainRemovedEventType    = domainEventPrefix + "removed"
)

func NewAddInstanceDomainUniqueConstraint(orgDomain string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueInstanceDomain,
		orgDomain,
		"Errors.Instance.Domain.AlreadyExists")
}

func NewRemoveInstanceDomainUniqueConstraint(orgDomain string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueInstanceDomain,
		orgDomain)
}

type DomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain    string `json:"domain,omitempty"`
	Generated bool   `json:"generated,omitempty"`
}

func (e *DomainAddedEvent) Data() interface{} {
	return e
}

func (e *DomainAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddInstanceDomainUniqueConstraint(e.Domain)}
}

func NewDomainAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string, generated bool) *DomainAddedEvent {
	return &DomainAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceDomainAddedEventType,
		),
		Domain:    domain,
		Generated: generated,
	}
}

func DomainAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	orgDomainAdded := &DomainAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INSTANCE-3noij", "unable to unmarshal instance domain added")
	}

	return orgDomainAdded, nil
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

func NewDomainPrimarySetEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainPrimarySetEvent {
	return &DomainPrimarySetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceDomainPrimarySetEventType,
		),
		Domain: domain,
	}
}

func DomainPrimarySetEventMapper(event *repository.Event) (eventstore.Event, error) {
	orgDomainAdded := &DomainPrimarySetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INSTANCE-29jöF", "unable to unmarshal instance domain added")
	}

	return orgDomainAdded, nil
}

type DomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain,omitempty"`
}

func (e *DomainRemovedEvent) Data() interface{} {
	return e
}

func (e *DomainRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveInstanceDomainUniqueConstraint(e.Domain)}
}

func NewDomainRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, domain string) *DomainRemovedEvent {
	return &DomainRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceDomainRemovedEventType,
		),
		Domain: domain,
	}
}

func DomainRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	orgDomainRemoved := &DomainRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgDomainRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INSTANCE-BngB2", "unable to unmarshal instance domain removed")
	}

	return orgDomainRemoved, nil
}
