package mirror

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type failedPayload struct {
	Cause string `json:"cause"`
}

type FailedEvent failedEvent
type failedEvent = eventstore.Event[failedPayload]

func FailedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*FailedEvent, error) {
	event, err := eventstore.EventFromStorage[failedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*FailedEvent)(event), nil
}

var (
	_ eventstore.Command = (*FailedCommand)(nil)
)

type FailedCommand struct {
	failedPayload
}

func NewFailedCommand(cause error) *FailedCommand {
	return &FailedCommand{
		failedPayload: failedPayload{
			Cause: cause.Error(),
		},
	}
}

// Creator implements eventstore.Command.
func (a *FailedCommand) Creator() string {
	return Creator
}

// Payload implements eventstore.Command.
func (a *FailedCommand) Payload() any {
	return a.failedPayload
}

// Revision implements [eventstore.Command].
func (*FailedCommand) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (e *FailedCommand) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// Type implements [eventstore.Command].
func (*FailedCommand) Type() string {
	return "system.mirror.failed"
}
