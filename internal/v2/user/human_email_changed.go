package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type humanEmailChangedPayload struct {
	Address domain.EmailAddress `json:"email,omitempty"`
}

type HumanEmailChangedEvent eventstore.Event[humanEmailChangedPayload]

const HumanEmailChangedType = humanPrefix + ".email.changed"

var _ eventstore.TypeChecker = (*HumanEmailChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanEmailChangedEvent) ActionType() string {
	return HumanEmailChangedType
}

func HumanEmailChangedEventFromStorage(event *eventstore.StorageEvent) (e *HumanEmailChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-Wr2lR", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanEmailChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanEmailChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
