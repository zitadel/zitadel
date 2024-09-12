package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgPrimaryDomain struct {
	Projection

	id string

	Domain string
}

func NewOrgPrimaryDomain(id string) *OrgPrimaryDomain {
	return &OrgPrimaryDomain{
		id: id,
	}
}

func (p *OrgPrimaryDomain) Reducers() Reducers {
	if p.Projection.Reducers != nil {
		return p.Projection.Reducers
	}

	p.Projection.Reducers = map[string]map[string]eventstore.ReduceEvent{
		org.AggregateType: {
			org.DomainPrimarySetType: p.reducePrimarySet,
		},
	}

	return p.Projection.Reducers
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
	p.Projection.Set(event)
	return nil
}

func (s *OrgPrimaryDomain) ShouldReduce(event *eventstore.StorageEvent) bool {
	return event.Aggregate.ID == s.id && s.Projection.ShouldReduce(event)
}
