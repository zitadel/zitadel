package readmodel

import (
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/system"
	"github.com/zitadel/zitadel/internal/v2/system/mirror"
)

type LastSuccessfulMirror struct {
	ID       string
	Position decimal.Decimal
	source   string
}

func NewLastSuccessfulMirror(source string) *LastSuccessfulMirror {
	return &LastSuccessfulMirror{
		source: source,
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
					mirror.SucceededType,
				),
				eventstore.EventCreatorsEqual(mirror.Creator),
			),
		),
		eventstore.FilterPagination(
			eventstore.Descending(),
			eventstore.Limit(1),
		),
	)
}

// Reduce implements eventstore.Reducer.
func (h *LastSuccessfulMirror) Reduce(events ...*eventstore.StorageEvent) (err error) {
	for _, event := range events {
		if event.Type == mirror.SucceededType {
			err = h.reduceSucceeded(event)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *LastSuccessfulMirror) reduceSucceeded(event *eventstore.StorageEvent) error {
	// if position is set we skip all older events
	if h.Position.GreaterThan(decimal.NewFromInt(0)) {
		return nil

	}
	succeededEvent, err := mirror.SucceededEventFromStorage(event)
	if err != nil {
		return err
	}

	if h.source != succeededEvent.Payload.Source {
		return nil
	}

	h.Position = succeededEvent.Payload.Position

	return nil
}
