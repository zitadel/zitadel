package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const SucceededType = eventTypePrefix + "succeeded"

type SucceededEvent eventstore.Event[eventstore.EmptyPayload]

var _ eventstore.TypeChecker = (*SucceededEvent)(nil)

func (e *SucceededEvent) ActionType() string {
	return SucceededType
}

func SucceededEventFromStorage(event *eventstore.StorageEvent) (e *SucceededEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "MIRRO-xh5IW", "Errors.Invalid.Event.Type")
	}

	return &SucceededEvent{
		StorageEvent: event,
	}, nil
}

func NewSucceededCommand() *eventstore.Command {
	return &eventstore.Command{
		Action: eventstore.Action[any]{
			Creator:  Creator,
			Type:     SucceededType,
			Revision: 1,
		},
	}
}
