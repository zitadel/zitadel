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
				eventstore.AggregateIDs(p.id),
				eventstore.AppendEvent(
					eventstore.SetEventTypes(
						org.AddedType,
						org.DeactivatedType,
						org.ReactivatedType,
						org.RemovedType,
					),
				),
			),
		),
	}
}

func (p *OrgState) Reduce(events ...*eventstore.StorageEvent) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}

		switch event.Type {
		case org.AddedType:
			p.State = org.ActiveState
		case org.DeactivatedType:
			p.State = org.InactiveState
		case org.ReactivatedType:
			p.State = org.ActiveState
		case org.RemovedType:
			p.State = org.RemovedState
		default:
			continue
		}
		p.position = event.Position
	}
	return nil
}
