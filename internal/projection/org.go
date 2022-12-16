package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

var _ Projection = (*Org)(nil)

type Org struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.OrgState
	Sequence      uint64

	Name   string
	Domain string

	instance string
}

func NewOrg(id, instance string) *Org {
	return &Org{
		ID:       id,
		instance: instance,
		State:    domain.OrgStateUnspecified,
	}
}

func (o *Org) Reduce(events []eventstore.Event) {
	for _, event := range events {
		o.ChangeDate = event.CreationDate()
		o.Sequence = event.Sequence()

		switch e := event.(type) {
		case *org.OrgAddedEvent:
			o.reduceAddedEvent(e)
		case *org.OrgChangedEvent:
			o.reduceChangedEvent(e)
		case *org.OrgDeactivatedEvent:
			o.reduceDeactivatedEvent(e)
		case *org.OrgReactivatedEvent:
			o.reduceReactivatedEvent(e)
		case *org.OrgRemovedEvent:
			o.reduceRemovedEvent(e)
		case *org.DomainPrimarySetEvent:
			o.reduceDomainPrimarySetEvent(e)
		}
	}
}

func (o *Org) SearchQuery(context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(o.instance).
		OrderAsc().
		ResourceOwner(o.ID).
		AddQuery().
		AggregateIDs(o.ID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDeactivatedEventType,
			org.OrgReactivatedEventType,
			org.OrgRemovedEventType,
			org.OrgDomainPrimarySetEventType,
		).
		Builder()
}

func (o *Org) reduceAddedEvent(event *org.OrgAddedEvent) {
	o.CreationDate = event.CreationDate()
	o.ResourceOwner = event.Aggregate().ResourceOwner
	o.State = domain.OrgStateActive
	o.Name = event.Name
}

func (o *Org) reduceChangedEvent(event *org.OrgChangedEvent) {
	o.Name = event.Name
}

func (o *Org) reduceDeactivatedEvent(event *org.OrgDeactivatedEvent) {
	o.State = domain.OrgStateInactive
}

func (o *Org) reduceReactivatedEvent(event *org.OrgReactivatedEvent) {
	o.State = domain.OrgStateActive
}

func (o *Org) reduceRemovedEvent(event *org.OrgRemovedEvent) {
	o.State = domain.OrgStateRemoved
}

func (o *Org) reduceDomainPrimarySetEvent(event *org.DomainPrimarySetEvent) {
	o.Domain = event.Domain
}
