package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type startedPayload struct {
	// Destination is the name of the database data are mirrored to
	Destination string `json:"destination"`
	// Either Instances or System needs to be set
	Instances []string `json:"instances,omitempty"`
	System    bool     `json:"system,omitempty"`
}

const StartedType = eventTypePrefix + "started"

type StartedEvent eventstore.Event[startedPayload]

var _ eventstore.TypeChecker = (*StartedEvent)(nil)

func (e *StartedEvent) ActionType() string {
	return StartedType
}

func StartedEventFromStorage(event *eventstore.StorageEvent) (e *StartedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "MIRRO-bwB9l", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[startedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &StartedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

func NewStartedSystemCommand(destination string) *eventstore.Command {
	return newStartedCommand(&startedPayload{
		Destination: destination,
		System:      true,
	})
}

func NewStartedInstancesCommand(destination string, instances []string) (*eventstore.Command, error) {
	if len(instances) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "MIRRO-8YkrE", "Errors.Mirror.NoInstances")
	}
	return newStartedCommand(&startedPayload{
		Destination: destination,
		Instances:   instances,
	}), nil
}

func newStartedCommand(payload *startedPayload) *eventstore.Command {
	return &eventstore.Command{
		Action: eventstore.Action[any]{
			Creator:  Creator,
			Type:     StartedType,
			Revision: 1,
			Payload:  *payload,
		},
	}
}
