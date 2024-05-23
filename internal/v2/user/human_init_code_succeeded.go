package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanInitCodeSucceededEvent eventstore.Event[eventstore.EmptyPayload]

const HumanInitCodeSucceededType = humanPrefix + ".initialization.check.succeeded"

var _ eventstore.TypeChecker = (*HumanInitCodeSucceededEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanInitCodeSucceededEvent) ActionType() string {
	return HumanInitCodeSucceededType
}

func HumanInitCodeSucceededEventFromStorage(event *eventstore.StorageEvent) (e *HumanInitCodeSucceededEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-12A5m", "Errors.Invalid.Event.Type")
	}

	return &HumanInitCodeSucceededEvent{
		StorageEvent: event,
	}, nil
}
