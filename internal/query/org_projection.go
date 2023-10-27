package query

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

var _ eventstore.QueryReducer = (*OrgProjection)(nil)

type OrgProjection struct {
	*eventstore.ReadModel
	Org
}

func newOrgProjection(instanceID, orgID string) *OrgProjection {
	return &OrgProjection{
		ReadModel: &eventstore.ReadModel{
			AggregateID: orgID,
			InstanceID:  instanceID,
		},
	}
}

// Query implements eventstore.QueryReducer.
func (p *OrgProjection) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(p.ReadModel.InstanceID).
		OrderAsc().
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(p.ReadModel.AggregateID).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDeactivatedEventType,
			org.OrgReactivatedEventType,
			org.OrgRemovedEventType,
			org.OrgDomainPrimarySetEventType,
		).
		Or().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(p.InstanceID).
		EventTypes(
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (p *OrgProjection) Reduce() error {
	for _, event := range p.Events {
		p.Org.ChangeDate = event.CreatedAt()
		p.Org.Sequence = event.Sequence()
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			p.Org.ID = event.Aggregate().ID
			p.Org.CreationDate = event.CreatedAt()
			p.Org.ResourceOwner = event.Aggregate().InstanceID
			p.Org.State = domain.OrgStateActive
			p.Org.Name = e.Name
		case *org.OrgChangedEvent:
			p.Org.Name = e.Name
		case *org.OrgDeactivatedEvent:
			p.Org.State = domain.OrgStateInactive
		case *org.OrgReactivatedEvent:
			p.Org.State = domain.OrgStateActive
		case *org.OrgRemovedEvent:
			p.Org.State = domain.OrgStateRemoved
		case *org.DomainPrimarySetEvent:
			p.Org.Domain = e.Domain
		case *instance.InstanceRemovedEvent:
			p.Org.State = domain.OrgStateRemoved
		}
	}
	return p.ReadModel.Reduce()
}
