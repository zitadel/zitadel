package query

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

var _ eventstore.QueryReducer = (*OrgStateProjection)(nil)

type OrgStateProjection struct {
	*eventstore.ReadModel
	state domain.OrgState
}

func (p *OrgStateProjection) Exists() bool {
	return p.state != domain.OrgStateRemoved && p.state != domain.OrgStateUnspecified
}

func newOrgStateProjection(instanceID, orgID string) *OrgStateProjection {
	return &OrgStateProjection{
		ReadModel: &eventstore.ReadModel{
			AggregateID: orgID,
			InstanceID:  instanceID,
		},
	}
}

// Query implements eventstore.QueryReducer.
func (p *OrgStateProjection) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(p.InstanceID).
		OrderAsc().
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(p.AggregateID).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgDeactivatedEventType,
			org.OrgReactivatedEventType,
			org.OrgRemovedEventType,
		).
		Or().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(p.InstanceID).
		EventTypes(
			instance.InstanceRemovedEventType,
		).
		Builder()
}

func (p *OrgStateProjection) Reduce() error {
	for _, event := range p.Events {
		switch event.Type() {
		case org.OrgAddedEventType:
			p.state = domain.OrgStateActive
		case org.OrgDeactivatedEventType:
			p.state = domain.OrgStateInactive
		case org.OrgReactivatedEventType:
			p.state = domain.OrgStateActive
		case org.OrgRemovedEventType, instance.InstanceRemovedEventType:
			p.state = domain.OrgStateRemoved
		}
	}
	return p.ReadModel.Reduce()
}
