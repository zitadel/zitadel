package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgState struct {
	projection

	id string

	State org.State
}

func NewStateProjection(ctx context.Context, id string) *OrgState {
	// TODO: check buffer for id and return from buffer if exists
	return &OrgState{
		projection: projection{
			instance: authz.GetInstance(ctx).InstanceID(),
		},
		id: id,
	}
}

func (p *OrgState) Filter() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(p.instance).
		PositionAfter(p.position).
		AddQuery().
		AggregateTypes(eventstore.AggregateType(org.AggregateType)).
		AggregateIDs(p.id).
		EventTypes(
			eventstore.EventType(org.Added.Type()),
			eventstore.EventType(org.Deactivated.Type()),
			eventstore.EventType(org.Reactivated.Type()),
			eventstore.EventType(org.Removed.Type()),
		).
		Builder()
}

func (p *OrgState) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		switch event.Type() {
		case eventstore.EventType(org.Added.Type()):
			p.State = org.ActiveState
		case eventstore.EventType(org.Deactivated.Type()):
			p.State = org.InactiveState
		case eventstore.EventType(org.Reactivated.Type()):
			p.State = org.ActiveState
		case eventstore.EventType(org.Removed.Type()):
			p.State = org.RemovedState
		default:
			continue
		}
		p.position = event.Position()
	}

	// TODO: if more than x events store state

	return nil
}
