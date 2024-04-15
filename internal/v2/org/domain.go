package org

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var uniqueOrgDomain = "org_domain"

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainAdded DomainAddedEvent
)

type DomainAddedEvent struct {
	*domain.AddedEvent
}

func DomainAddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainAddedEvent, error) {
	event, err := domain.AddedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainAddedEvent{
		AddedEvent: event,
	}, nil
}

func (e DomainAddedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "org") && e.AddedEvent.HasTypeSuffix(typ)
}

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainVerified DomainVerifiedEvent
)

type DomainVerifiedEvent struct {
	*domain.VerifiedEvent
}

func DomainVerifiedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainVerifiedEvent, error) {
	event, err := domain.VerifiedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainVerifiedEvent{
		VerifiedEvent: event,
	}, nil
}

func (e DomainVerifiedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "org") && e.VerifiedEvent.HasTypeSuffix(typ)
}

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainPrimarySet DomainPrimarySetEvent
)

type DomainPrimarySetEvent struct {
	*domain.PrimarySetEvent
}

func DomainPrimarySetEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPrimarySetEvent, error) {
	event, err := domain.PrimarySetEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainPrimarySetEvent{
		PrimarySetEvent: event,
	}, nil
}

func (e DomainPrimarySetEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "org") && e.PrimarySetEvent.HasTypeSuffix(typ)
}

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainRemoved DomainRemovedEvent
)

type DomainRemovedEvent struct {
	*domain.RemovedEvent
}

func DomainRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainRemovedEvent, error) {
	event, err := domain.RemovedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainRemovedEvent{
		RemovedEvent: event,
	}, nil
}

func (e DomainRemovedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "org") && e.RemovedEvent.HasTypeSuffix(typ)
}
