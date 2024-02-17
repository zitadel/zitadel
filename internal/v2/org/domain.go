package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var uniqueOrgDomain = "org_domain"

var (
	_           eventstore.Command = (*DomainAddedEvent)(nil)
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

func (e *DomainAddedEvent) Type() string {
	return string(append([]byte("org."), e.AddedEvent.Type()...))
}

// UniqueConstraints implements eventstore.Command.
func (e *DomainAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

var (
	_              eventstore.Command = (*DomainVerifiedEvent)(nil)
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

func (e *DomainVerifiedEvent) Type() string {
	return string(append([]byte("org."), e.VerifiedEvent.Type()...))
}

// UniqueConstraints implements eventstore.Command.
func (e *DomainVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(uniqueOrgDomain, e.Name, "Errors.Org.Domain.AlreadyExists"),
	}
}

var (
	_                eventstore.Command = (*SetDomainPrimaryEvent)(nil)
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

func (e *SetDomainPrimaryEvent) Type() string {
	return string(append([]byte("org."), e.PrimarySetEvent.Type()...))
}

// UniqueConstraints implements eventstore.Command.
func (*SetDomainPrimaryEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
