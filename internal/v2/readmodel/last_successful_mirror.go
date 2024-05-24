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
					mirror.StartedType,
					mirror.SucceededType,
				),
				eventstore.EventCreatorsEqual(mirror.Creator),
			),
		),
	)
}

// Reduce implements eventstore.Reducer.
func (h *LastSuccessfulMirror) Reduce(events ...*eventstore.StorageEvent) (err error) {
	for _, event := range events {
		switch event.Type {
		case mirror.StartedType:
			err = h.reduceStarted(event)
		case mirror.SucceededType:
			err = h.reduceSucceeded(event)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *LastSuccessfulMirror) reduceStarted(event *eventstore.StorageEvent) error {
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

func (h *LastSuccessfulMirror) reduceSucceeded(event *eventstore.StorageEvent) error {
	if h.ID != event.Aggregate.ID {
		return nil
	}

	h.Position = event.Position

	return nil
}
