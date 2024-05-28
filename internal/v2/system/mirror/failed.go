package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type failedPayload struct {
	Cause string `json:"cause"`
	// Source is the name of the database data are mirrored to
	Source string `json:"source"`
}

const FailedType = eventTypePrefix + "failed"

type FailedEvent eventstore.Event[failedPayload]

var _ eventstore.TypeChecker = (*FailedEvent)(nil)

func (e *FailedEvent) ActionType() string {
	return FailedType
}

func FailedEventFromStorage(event *eventstore.StorageEvent) (e *FailedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "MIRRO-bwB9l", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[failedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &FailedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

func NewFailedCommand(source string, cause error) *eventstore.Command {
	return &eventstore.Command{
		Action: eventstore.Action[any]{
			Creator: Creator,
			Type:    FailedType,
			Payload: failedPayload{
				Cause:  cause.Error(),
				Source: source,
			},
			Revision: 1,
		},
	}
}
