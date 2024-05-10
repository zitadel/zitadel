package readmodel

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/system"
	"github.com/zitadel/zitadel/internal/v2/system/mirror"
)

type LastSuccessfulMirror struct {
	ID          string
	Position    eventstore.GlobalPosition
	destination string
}

func NewLastSuccessfulMirror(destination string) *LastSuccessfulMirror {
	return &LastSuccessfulMirror{
		destination: destination,
	}
}

var _ eventstore.Reducer = (*LastSuccessfulMirror)(nil)

func (p *LastSuccessfulMirror) Filter() *eventstore.Filter {
	return eventstore.NewFilter(
		eventstore.AppendAggregateFilter(
			system.AggregateType,
			eventstore.AggregateOwnersEqual(system.AggregateOwner),
			eventstore.AppendEvent(
				eventstore.SetEventTypes(
					"system.mirror.started",
					"system.mirror.succeeded",
				),
				eventstore.EventCreatorsEqual(mirror.Creator),
			),
		),
	)
}

// Reduce implements eventstore.Reducer.
func (h *LastSuccessfulMirror) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) (err error) {
	for _, event := range events {
		switch event.Type {
		case "system.mirror.started":
			err = h.reduceStarted(event)
		case "system.mirror.succeeded":
			err = h.reduceSucceeded(event)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *LastSuccessfulMirror) reduceStarted(event *eventstore.Event[eventstore.StoragePayload]) error {
	startedEvent, err := mirror.StartedEventFromStorage(event)
	if err != nil {
		return err
	}

	if h.destination != startedEvent.Payload.Destination {
		return nil
	}

	h.ID = event.Aggregate.ID

	return nil
}

func (h *LastSuccessfulMirror) reduceSucceeded(event *eventstore.Event[eventstore.StoragePayload]) error {
	if h.ID != event.Aggregate.ID {
		return nil
	}

	h.Position = event.Position

	return nil
}
