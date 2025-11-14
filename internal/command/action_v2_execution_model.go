package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/execution"
)

type ExecutionWriteModel struct {
	eventstore.WriteModel

	Targets          []string
	Includes         []string
	ExecutionTargets []*execution.Target
}

func (e *ExecutionWriteModel) ExecutionTargetsEqual(targets []*execution.Target) bool {
	if len(e.ExecutionTargets) != len(targets) {
		return false
	}
	for i := range e.ExecutionTargets {
		if e.ExecutionTargets[i].Type != targets[i].Type || e.ExecutionTargets[i].Target != targets[i].Target {
			return false
		}
	}
	return true
}

func (e *ExecutionWriteModel) IncludeList() []string {
	includes := make([]string, 0)
	for i := range e.ExecutionTargets {
		if e.ExecutionTargets[i].Type == domain.ExecutionTargetTypeInclude {
			includes = append(includes, e.ExecutionTargets[i].Target)
		}
	}
	return includes
}

func (e *ExecutionWriteModel) Exists() bool {
	return len(e.ExecutionTargets) > 0 || len(e.Includes) > 0 || len(e.Targets) > 0
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
		case *execution.SetEventV2:
			wm.ExecutionTargets = e.Targets
		case *execution.RemovedEvent:
			wm.Targets = nil
			wm.Includes = nil
			wm.ExecutionTargets = nil
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
			execution.SetEventV2Type,
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

type ExecutionsExistWriteModel struct {
	eventstore.WriteModel

	ids         []string
	existingIDs []string
}

func (e *ExecutionsExistWriteModel) AllExists() bool {
	return len(e.ids) == len(e.existingIDs)
}

func NewExecutionsExistWriteModel(ids []string, resourceOwner string) *ExecutionsExistWriteModel {
	return &ExecutionsExistWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
		ids: ids,
	}
}

func (wm *ExecutionsExistWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *execution.SetEvent:
			if !slices.Contains(wm.existingIDs, e.Aggregate().ID) {
				wm.existingIDs = append(wm.existingIDs, e.Aggregate().ID)
			}
		case *execution.SetEventV2:
			if !slices.Contains(wm.existingIDs, e.Aggregate().ID) {
				wm.existingIDs = append(wm.existingIDs, e.Aggregate().ID)
			}
		case *execution.RemovedEvent:
			i := slices.Index(wm.existingIDs, e.Aggregate().ID)
			if i >= 0 {
				wm.existingIDs = slices.Delete(wm.existingIDs, i, i+1)
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ExecutionsExistWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(execution.AggregateType).
		AggregateIDs(wm.ids...).
		EventTypes(execution.SetEventType,
			execution.SetEventV2Type,
			execution.RemovedEventType).
		Builder()
}
