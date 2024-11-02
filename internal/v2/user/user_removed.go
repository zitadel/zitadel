package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type RemovedEvent eventstore.Event[eventstore.EmptyPayload]

const RemovedType = AggregateType + ".removed"

var _ eventstore.TypeChecker = (*RemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *RemovedEvent) ActionType() string {
	return RemovedType
}

func RemovedEventFromStorage(event *eventstore.StorageEvent) (e *RemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-UN6Xa", "Errors.Invalid.Event.Type")
	}

	return &RemovedEvent{
		StorageEvent: event,
	}, nil
}
