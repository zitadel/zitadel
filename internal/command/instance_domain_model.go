package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type InstanceDomainWriteModel struct {
	eventstore.WriteModel

	Domain string
	State  domain.InstanceDomainState
}

func NewInstanceDomainWriteModel(instanceID string, instanceDomain string) *InstanceDomainWriteModel {
	return &InstanceDomainWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
		Domain: instanceDomain,
	}
}

func (wm *InstanceDomainWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.DomainAddedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *iam.DomainRemovedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceDomainWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.DomainAddedEvent:
			wm.Domain = e.Domain
			wm.State = domain.InstanceDomainStateActive
		case *iam.DomainRemovedEvent:
			wm.State = domain.InstanceDomainStateRemoved
		}
	}
	return nil
}

func (wm *InstanceDomainWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.InstanceDomainAddedEventType,
			iam.InstanceDomainRemovedEventType).
		Builder()
}

type InstanceDomainsWriteModel struct {
	eventstore.WriteModel

	Domains []*InstanceDomain
}

type InstanceDomain struct {
	Domain string
	State  domain.InstanceDomainState
}

func NewInstanceDomainsWriteModel(orgID string) *InstanceDomainsWriteModel {
	return &InstanceDomainsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
		Domains: make([]*InstanceDomain, 0),
	}
}

func (wm *InstanceDomainsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.DomainAddedEvent:
			wm.Domains = append(wm.Domains, &InstanceDomain{Domain: e.Domain, State: domain.InstanceDomainStateActive})
		case *iam.DomainRemovedEvent:
			for _, d := range wm.Domains {
				if d.Domain == e.Domain {
					d.State = domain.InstanceDomainStateRemoved
					continue
				}
			}
		}
	}
	return nil
}

func (wm *InstanceDomainsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.InstanceDomainAddedEventType,
			iam.InstanceDomainRemovedEventType).
		Builder()
}
