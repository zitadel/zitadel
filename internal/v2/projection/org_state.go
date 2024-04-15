package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgState struct {
	projection

	id string

	org.State
}

func NewStateProjection(id string) *OrgState {
	// TODO: check buffer for id and return from buffer if exists
	return &OrgState{
		id: id,
	}
}

func (p *OrgState) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.FilterPagination(
				eventstore.Descending(),
				eventstore.GlobalPositionGreater(&p.position),
			),
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.AggregateID(p.id),
				eventstore.AppendEvent(
					eventstore.EventTypes(
						"org.added",
						"org.deactivated",
						"org.reactivated",
						"org.removed",
					),
				),
			),
		),
	}
}

func (p *OrgState) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}

		switch {
		case org.Added.IsType(event.Type):
			p.State = org.ActiveState
		case org.Deactivated.IsType(event.Type):
			p.State = org.InactiveState
		case org.Reactivated.IsType(event.Type):
			p.State = org.ActiveState
		case org.Removed.IsType(event.Type):
			p.State = org.RemovedState
		default:
			continue
		}
		p.position = event.Position
	}

	// TODO: if more than x events store state

	return nil
}
