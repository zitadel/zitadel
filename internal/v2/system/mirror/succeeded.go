package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type SucceededEvent succeededEvent
type succeededEvent = eventstore.Event[struct{}]

func SucceededEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*SucceededEvent, error) {
	event, err := eventstore.EventFromStorage[succeededEvent](e)
	if err != nil {
		return nil, err
	}
	return (*SucceededEvent)(event), nil
}

var (
	_ eventstore.Command = (*SucceededCommand)(nil)
)

type SucceededCommand struct {
}

func NewSucceededCommand() *SucceededCommand {
	return new(SucceededCommand)
}

// Creator implements eventstore.Command.
func (a *SucceededCommand) Creator() string {
	return Creator
}

// Payload implements eventstore.Command.
func (a *SucceededCommand) Payload() any {
	return nil
}

// Revision implements [eventstore.Command].
func (*SucceededCommand) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *SucceededCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// Type implements [eventstore.Command].
func (*SucceededCommand) Type() string {
	return "system.mirror.succeeded"
}
