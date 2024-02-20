package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var uniqueOrgDomain = "org_domain"

var (
	_ eventstore.Command = (*DomainAddedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainAdded *DomainAddedEvent
)

type DomainAddedEvent struct {
	*domain.AddedEvent
}

func NewDomainAddedEvent(ctx context.Context, name string) (*DomainAddedEvent, error) {
	event, err := domain.NewAddedEvent(ctx, name)
	if err != nil {
		return nil, err
	}
	return &DomainAddedEvent{AddedEvent: event}, nil
}

// Type implements [eventstore.action].
func (e *DomainAddedEvent) Type() string {
	return string(append([]byte(eventTypePrefix), new(domain.AddedEvent).Type()...))
}

// UniqueConstraints implements [eventstore.Command].
func (e *DomainAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

var (
	_ eventstore.Command = (*DomainVerifiedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainVerified *DomainVerifiedEvent
)

type DomainVerifiedEvent struct {
	*domain.VerifiedEvent
}

func NewDomainVerifiedEvent(ctx context.Context, name string) (*DomainVerifiedEvent, error) {
	event, err := domain.NewVerifiedEvent(ctx, name)
	if err != nil {
		return nil, err
	}
	return &DomainVerifiedEvent{VerifiedEvent: event}, nil
}

// Type implements [eventstore.action].
func (e *DomainVerifiedEvent) Type() string {
	return string(append([]byte(eventTypePrefix), new(domain.VerifiedEvent).Type()...))
}

// UniqueConstraints implements [eventstore.Command].
func (e *DomainVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(uniqueOrgDomain, e.Name, "Errors.Org.Domain.AlreadyExists"),
	}
}

var (
	_ eventstore.Command = (*SetDomainPrimaryEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainSetPrimary *SetDomainPrimaryEvent
)

type SetDomainPrimaryEvent struct {
	*domain.PrimarySetEvent
}

func NewSetDomainPrimaryEvent(ctx context.Context, name string) (*SetDomainPrimaryEvent, error) {
	event, err := domain.NewSetPrimaryEvent(ctx, name)
	if err != nil {
		return nil, err
	}
	return &SetDomainPrimaryEvent{PrimarySetEvent: event}, nil
}

// Type implements [eventstore.action].
func (e *SetDomainPrimaryEvent) Type() string {
	return string(append([]byte(eventTypePrefix), new(domain.PrimarySetEvent).Type()...))
}

// UniqueConstraints implements [eventstore.Command].
func (*SetDomainPrimaryEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
