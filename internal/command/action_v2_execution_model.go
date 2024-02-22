package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/execution"
)

type ExecutionWriteModel struct {
	eventstore.WriteModel

	Targets  []string
	Includes []string
}

func (e *ExecutionWriteModel) Exists() bool {
	return len(e.Targets) > 0 || len(e.Includes) > 0
}

func NewExecutionWriteModel(id string, resourceOwner string) *ExecutionWriteModel {
	return &ExecutionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
	}
}

func (wm *ExecutionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *execution.SetEvent:
			wm.Targets = e.Targets
			wm.Includes = e.Includes
		case *execution.RemovedEvent:
			wm.Targets = nil
			wm.Includes = nil
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ExecutionWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(execution.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(execution.SetEventType,
			execution.RemovedEventType).
		Builder()
}

func ExecutionAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          execution.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       execution.AggregateVersion,
	}
}
