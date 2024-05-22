package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DomainPolicyAddedType = eventTypePrefix + policy.DomainPolicyAddedTypeSuffix

type DomainPolicyAddedPayload policy.DomainPolicyAddedPayload

type DomainPolicyAddedEvent eventstore.Event[DomainPolicyAddedPayload]

var _ eventstore.TypeChecker = (*DomainPolicyAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainPolicyAddedEvent) ActionType() string {
	return DomainPolicyAddedType
}

func DomainPolicyAddedEventFromStorage(event *eventstore.StorageEvent) (e *DomainPolicyAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-z1a7D", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainPolicyAddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainPolicyAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

const DomainPolicyChangedType = eventTypePrefix + policy.DomainPolicyChangedTypeSuffix

type DomainPolicyChangedPayload policy.DomainPolicyChangedPayload

type DomainPolicyChangedEvent eventstore.Event[DomainPolicyChangedPayload]

var _ eventstore.TypeChecker = (*DomainPolicyChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainPolicyChangedEvent) ActionType() string {
	return DomainPolicyChangedType
}

func DomainPolicyChangedEventFromStorage(event *eventstore.StorageEvent) (e *DomainPolicyChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-BTLhd", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[DomainPolicyChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainPolicyChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
