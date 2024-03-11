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
			eventstore.Descending(),
			eventstore.WithPositionAtLeast(p.position),
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.WithAggregateID(p.id),
				eventstore.AppendEvent(
					eventstore.WithEventType(org.Added.Type()),
				),
				eventstore.AppendEvent(
					eventstore.WithEventType(org.Deactivated.Type()),
				),
				eventstore.AppendEvent(
					eventstore.WithEventType(org.Reactivated.Type()),
				),
				eventstore.AppendEvent(
					eventstore.WithEventType(org.Removed.Type()),
				),
			),
		),
	}
}

func (p *OrgState) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}

		switch event.Type() {
		case org.Added.Type():
			p.State = org.ActiveState
		case org.Deactivated.Type():
			p.State = org.InactiveState
		case org.Reactivated.Type():
			p.State = org.ActiveState
		case org.Removed.Type():
			p.State = org.RemovedState
		default:
			continue
		}
		p.position = event.Position()
	}

	// TODO: if more than x events store state

	return nil
}
