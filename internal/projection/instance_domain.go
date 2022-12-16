package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var _ Projection = (*SearchInstanceDomain)(nil)

type SearchInstanceDomain struct {
	host             string
	removedInstances []string

	InstanceID string
}

func NewSearchInstanceDomain(host string) *SearchInstanceDomain {
	return &SearchInstanceDomain{
		host: host,
	}
}

func (domains *SearchInstanceDomain) Reduce(events []eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DomainAddedEvent:
			if ok := domains.reduceAdded(e); ok {
				return
			}
		case *instance.DomainRemovedEvent:
			domains.reduceRemoved(e)
		case *instance.InstanceRemovedEvent:
			domains.reduceInstanceRemoved(e)
		}
	}
}

func (domains *SearchInstanceDomain) SearchQuery(context.Context) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
		).
		EventData(map[string]interface{}{
			"domain": domains.host,
		}).
		Builder()
}

func (domains *SearchInstanceDomain) reduceAdded(event *instance.DomainAddedEvent) (ok bool) {
	for _, removed := range domains.removedInstances {
		if removed == event.Domain {
			return false
		}
	}
	domains.InstanceID = event.Aggregate().ID
	return true
}

func (domains *SearchInstanceDomain) reduceRemoved(event *instance.DomainRemovedEvent) {
	domains.removedInstances = append(domains.removedInstances, event.Aggregate().ID)
}

func (domains *SearchInstanceDomain) reduceInstanceRemoved(event *instance.InstanceRemovedEvent) {
	domains.removedInstances = append(domains.removedInstances, event.Aggregate().ID)
}
