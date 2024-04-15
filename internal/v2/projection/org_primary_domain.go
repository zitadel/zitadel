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
			eventstore.FilterPagination(
				eventstore.GlobalPositionGreater(&p.position),
			),
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.AggregateID(p.id),
				eventstore.AppendEvent(
					eventstore.EventType("org.domain.primary.set"),
				),
			),
		),
	}
}

func (p *OrgPrimaryDomain) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}

		if !org.DomainPrimarySet.IsType(event.Type) {
			continue
		}
		e, err := org.DomainPrimarySetEventFromStorage(event)
		if err != nil {
			return err
		}

		p.Domain = e.Payload.Name
		p.projection.reduce(event)
	}

	return nil
}
