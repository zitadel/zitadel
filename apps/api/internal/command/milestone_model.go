package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

type MilestonesReachedWriteModel struct {
	eventstore.WriteModel
	MilestonesReached
}

func NewMilestonesReachedWriteModel(instanceID string) *MilestonesReachedWriteModel {
	return &MilestonesReachedWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: instanceID,
			InstanceID:  instanceID,
		},
		MilestonesReached: MilestonesReached{
			InstanceID: instanceID,
		},
	}
}

func (m *MilestonesReachedWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(milestone.AggregateType).
		AggregateIDs(m.AggregateID).
		EventTypes(milestone.ReachedEventType, milestone.PushedEventType).
		Builder()
}

func (m *MilestonesReachedWriteModel) Reduce() error {
	for _, event := range m.Events {
		if e, ok := event.(*milestone.ReachedEvent); ok {
			m.reduceReachedEvent(e)
		}
	}
	return m.WriteModel.Reduce()
}

func (m *MilestonesReachedWriteModel) reduceReachedEvent(e *milestone.ReachedEvent) {
	switch e.MilestoneType {
	case milestone.InstanceCreated:
		m.InstanceCreated = true
	case milestone.AuthenticationSucceededOnInstance:
		m.AuthenticationSucceededOnInstance = true
	case milestone.ProjectCreated:
		m.ProjectCreated = true
	case milestone.ApplicationCreated:
		m.ApplicationCreated = true
	case milestone.AuthenticationSucceededOnApplication:
		m.AuthenticationSucceededOnApplication = true
	case milestone.InstanceDeleted:
		m.InstanceDeleted = true
	}
}
