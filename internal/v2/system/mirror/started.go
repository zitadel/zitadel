package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type startedPayload struct {
	Destination string `json:"destination"`
	// Either Instances or System needs to be set
	Instances []string `json:"instances,omitempty"`
	System    bool     `json:"system,omitempty"`
}

type StartedEvent startedEvent
type startedEvent = eventstore.Event[startedPayload]

func StartedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*StartedEvent, error) {
	event, err := eventstore.EventFromStorage[startedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*StartedEvent)(event), nil
}

var (
	_ eventstore.Command = (*StartedCommand)(nil)
)

type StartedCommand struct {
	startedPayload
}

func NewStartedInstancesCommand(destination string, instances []string) *StartedCommand {
	return &StartedCommand{
		startedPayload: startedPayload{
			Destination: destination,
			Instances:   instances,
		},
	}
}

func NewStartedSystemCommand(destination string) *StartedCommand {
	return &StartedCommand{
		startedPayload: startedPayload{
			Destination: destination,
			System:      true,
		},
	}
}

// Creator implements eventstore.Command.
func (a *StartedCommand) Creator() string {
	return Creator
}

// Payload implements eventstore.Command.
func (a *StartedCommand) Payload() any {
	return a.startedPayload
}

// Revision implements [eventstore.Command].
func (*StartedCommand) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *StartedCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// Type implements [eventstore.Command].
func (*StartedCommand) Type() string {
	return "system.mirror.started"
}
