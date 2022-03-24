package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type InstanceWriteModel struct {
	eventstore.WriteModel

	Name            string
	State           domain.InstanceState
	GeneratedDomain string
}

func NewInstanceWriteModel(instanceID string) *InstanceWriteModel {
	return &InstanceWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
	}
}

func (wm *InstanceWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.InstanceAddedEvent:
			wm.Name = e.Name
			wm.State = domain.InstanceStateActive
		case *iam.InstanceChangedEvent:
			wm.Name = e.Name
		case *iam.InstanceRemovedEvent:
			wm.State = domain.InstanceStateRemoved
		}
	}
	return nil
}

func (wm *InstanceWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.InstanceAddedEventType,
			iam.InstanceChangedEventType,
			iam.InstanceRemovedEventType).
		Builder()
}

func InstanceAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion)
}
