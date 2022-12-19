package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var _ Projection = (*SearchInstanceDomain)(nil)

type SearchInstanceDomain struct {
	domain string
	// removed key is instance id
	removed map[string]*removedInstance

	InstanceID string
}

type removedInstance struct {
	isRemoved bool
	// domains are the list of removed domains
	domains []string
}

func NewSearchInstanceDomain(domain string) *SearchInstanceDomain {
	return &SearchInstanceDomain{
		domain:  domain,
		removed: map[string]*removedInstance{},
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
			domains.reduceDomainRemoved(e)
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
			"domain": domains.domain,
		}).
		Or().
		AggregateTypes(instance.AggregateType).
		EventTypes(instance.InstanceRemovedEventType).
		Builder()
}

func (domains *SearchInstanceDomain) reduceAdded(event *instance.DomainAddedEvent) (ok bool) {
	if domains.isRemoved(event.Aggregate().ID, event.Domain) {
		return false
	}
	domains.InstanceID = event.Aggregate().ID
	return true
}

func (domains *SearchInstanceDomain) reduceDomainRemoved(event *instance.DomainRemovedEvent) {
	if _, ok := domains.removed[event.Aggregate().ID]; !ok {
		domains.removed[event.Aggregate().ID] = new(removedInstance)
	}
	domains.removed[event.Aggregate().ID].domains =
		append(domains.removed[event.Aggregate().ID].domains, event.Domain)
}

func (domains *SearchInstanceDomain) reduceInstanceRemoved(event *instance.InstanceRemovedEvent) {
	if _, ok := domains.removed[event.Aggregate().ID]; !ok {
		domains.removed[event.Aggregate().ID] = new(removedInstance)
	}
	domains.removed[event.Aggregate().ID].isRemoved = true
}

func (domains *SearchInstanceDomain) isRemoved(instanceID, domain string) bool {
	removed, ok := domains.removed[instanceID]
	if !ok {
		return false
	}

	if removed.isRemoved {
		return true
	}

	for _, removedDomain := range removed.domains {
		if removedDomain == domain {
			return true
		}
	}
	return false
}
