package org

import (
	"github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DomainAddedType = "org." + domain.AddedTypeSuffix

type DomainAddedPayload domain.AddedPayload

type DomainAddedEvent eventstore.Event[DomainAddedPayload]

var _ eventstore.TypeChecker = (*DomainAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainAddedEvent) ActionType() string {
	return DomainAddedType
}

func DomainAddedEventFromStorage(event *eventstore.StorageEvent) (e *DomainAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-CXVe3", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainAddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainVerifiedType = "org." + domain.VerifiedTypeSuffix

type DomainVerifiedPayload domain.VerifiedPayload

type DomainVerifiedEvent eventstore.Event[DomainVerifiedPayload]

var _ eventstore.TypeChecker = (*DomainVerifiedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainVerifiedEvent) ActionType() string {
	return DomainVerifiedType
}

func DomainVerifiedEventFromStorage(event *eventstore.StorageEvent) (e *DomainVerifiedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-RAwdb", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainVerifiedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainVerifiedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainPrimarySetType = "org." + domain.PrimarySetTypeSuffix

type DomainPrimarySetPayload domain.PrimarySetPayload

type DomainPrimarySetEvent eventstore.Event[DomainPrimarySetPayload]

var _ eventstore.TypeChecker = (*DomainPrimarySetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainPrimarySetEvent) ActionType() string {
	return DomainPrimarySetType
}

func DomainPrimarySetEventFromStorage(event *eventstore.StorageEvent) (e *DomainPrimarySetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-7P3Iz", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainPrimarySetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainPrimarySetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainRemovedType = "org." + domain.RemovedTypeSuffix

type DomainRemovedPayload domain.RemovedPayload

type DomainRemovedEvent eventstore.Event[DomainRemovedPayload]

var _ eventstore.TypeChecker = (*DomainRemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainRemovedEvent) ActionType() string {
	return DomainRemovedType
}

func DomainRemovedEventFromStorage(event *eventstore.StorageEvent) (e *DomainRemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-ndpL2", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainRemovedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainRemovedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
