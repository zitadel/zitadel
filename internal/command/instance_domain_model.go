package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceDomainWriteModel struct {
	eventstore.WriteModel

	Domain    string
	Generated bool
	State     domain.InstanceDomainState
}

func NewInstanceDomainWriteModel(instanceID, instanceDomain string) *InstanceDomainWriteModel {
	return &InstanceDomainWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
		Domain: instanceDomain,
	}
}

func (wm *InstanceDomainWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DomainAddedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *instance.DomainRemovedEvent:
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
		case *instance.DomainAddedEvent:
			wm.Domain = e.Domain
			wm.Generated = e.Generated
			wm.State = domain.InstanceDomainStateActive
		case *instance.DomainRemovedEvent:
			wm.State = domain.InstanceDomainStateRemoved
		}
	}
	return nil
}

func (wm *InstanceDomainWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType).
		Builder()
}
