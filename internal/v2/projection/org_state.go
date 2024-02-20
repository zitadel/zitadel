package projection

import (
	"context"

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

func (p *OrgState) Filter(ctx context.Context) *eventstore.Filter {
	return eventstore.NewFilter(
		ctx,
		eventstore.FilterPositionAtLeast(p.position),
		eventstore.FilterEventQuery(
			eventstore.FilterAggregateTypes(org.AggregateType),
			eventstore.FilterAggregateIDs(p.id),
			eventstore.FilterEventTypes(
				org.Added.Type(),
				org.Deactivated.Type(),
				org.Reactivated.Type(),
				org.Removed.Type(),
			),
		),
	)
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
