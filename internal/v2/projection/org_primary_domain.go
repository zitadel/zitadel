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

func (p *OrgPrimaryDomain) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.WithPositionAtLeast(p.position, 0),
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.WithAggregateID(p.id),
				eventstore.AppendEvent(
					eventstore.WithEventType(org.DomainSetPrimary.Type()),
				),
			),
		),
	}
}

func (p *OrgPrimaryDomain) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}
		if event.Type() != org.DomainSetPrimary.Type() {
			continue
		}

		e := new(org.SetDomainPrimaryEvent)
		if err := event.Unmarshal(e); err != nil {
			return err
		}

		p.Domain = e.Name
		p.position = event.Position()
		p.sequence = event.Sequence()
	}

	return nil
}
