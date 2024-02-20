package projection

import (
	"context"

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

func (p *OrgPrimaryDomain) Filter(ctx context.Context) *eventstore.Filter {
	return eventstore.NewFilter(
		ctx,
		eventstore.FilterPositionAtLeast(p.position),
		eventstore.FilterEventQuery(
			eventstore.FilterAggregateTypes(org.AggregateType),
			eventstore.FilterAggregateIDs(p.id),
			eventstore.FilterEventTypes(org.DomainSetPrimary.Type()),
		),
	)
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
