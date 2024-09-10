package instance

import (
	"github.com/zitadel/zitadel/internal/v2/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DomainAddedType = eventTypePrefix + domain.AddedTypeSuffix

type DomainAddedEvent eventstore.Event[domain.AddedPayload]

var _ eventstore.TypeChecker = (*DomainAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainAddedEvent) ActionType() string {
	return DomainAddedType
}

func DomainAddedEventFromStorage(event *eventstore.StorageEvent) (e *DomainAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-CXVe3", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[domain.AddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainVerifiedType = "org." + domain.VerifiedTypeSuffix

type DomainVerifiedEvent eventstore.Event[domain.VerifiedPayload]

var _ eventstore.TypeChecker = (*DomainVerifiedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainVerifiedEvent) ActionType() string {
	return DomainVerifiedType
}

func DomainVerifiedEventFromStorage(event *eventstore.StorageEvent) (e *DomainVerifiedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-RAwdb", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[domain.VerifiedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainVerifiedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainPrimarySetType = "org." + domain.PrimarySetTypeSuffix

type DomainPrimarySetEvent eventstore.Event[domain.PrimarySetPayload]

var _ eventstore.TypeChecker = (*DomainPrimarySetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainPrimarySetEvent) ActionType() string {
	return DomainPrimarySetType
}

func DomainPrimarySetEventFromStorage(event *eventstore.StorageEvent) (e *DomainPrimarySetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-7P3Iz", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[domain.PrimarySetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainPrimarySetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainRemovedType = "org." + domain.RemovedTypeSuffix

type DomainRemovedEvent eventstore.Event[domain.RemovedPayload]

var _ eventstore.TypeChecker = (*DomainRemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainRemovedEvent) ActionType() string {
	return DomainRemovedType
}

func DomainRemovedEventFromStorage(event *eventstore.StorageEvent) (e *DomainRemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-ndpL2", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[domain.RemovedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainRemovedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
