package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgPrimaryDomain struct {
	projection

	id string

	Domain string
}

func NewOrgPrimaryDomain(id string) *OrgPrimaryDomain {
	return &OrgPrimaryDomain{
		id: id,
	}
}

// func (p *OrgPrimaryDomain) Filter() []*eventstore.Filter {
// 	return []*eventstore.Filter{
// 		eventstore.NewFilter(
// 			eventstore.FilterPagination(
// 				eventstore.GlobalPositionGreater(&p.position),
// 			),
// 			eventstore.AppendAggregateFilter(
// 				org.AggregateType,
// 				eventstore.AggregateIDs(p.id),
// 				eventstore.AppendEvent(
// 					eventstore.SetEventTypes(org.DomainPrimarySetType),
// 				),
// 			),
// 		),
// 	}
// }

func (p *OrgPrimaryDomain) Reducers() map[string]map[string]eventstore.ReduceEvent {
	if p.reducers != nil {
		return p.reducers
	}

	p.reducers = map[string]map[string]eventstore.ReduceEvent{
		org.AggregateType: {
			org.DomainPrimarySetType: p.reducePrimarySet,
		},
	}

	return p.reducers
}

func (p *OrgPrimaryDomain) reducePrimarySet(event *eventstore.StorageEvent) error {
	if !p.ShouldReduce(event) {
		return nil
	}

	e, err := org.DomainPrimarySetEventFromStorage(event)
	if err != nil {
		return err
	}

	p.Domain = e.Payload.Name
	p.projection.set(event)
	return nil
}

func (s *OrgPrimaryDomain) ShouldReduce(event *eventstore.StorageEvent) bool {
	return event.Aggregate.ID == s.id && s.projection.ShouldReduce(event)
}
