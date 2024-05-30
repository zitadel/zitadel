package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type domainClaimedPayload struct {
	Username          string `json:"userName"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
}

type DomainClaimedEvent eventstore.Event[domainClaimedPayload]

const DomainClaimedType = AggregateType + ".domain.claimed.sent"

var _ eventstore.TypeChecker = (*DomainClaimedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DomainClaimedEvent) ActionType() string {
	return DomainClaimedType
}

func DomainClaimedEventFromStorage(event *eventstore.StorageEvent) (e *DomainClaimedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-x8O4o", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[domainClaimedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DomainClaimedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
