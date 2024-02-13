package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/execution"
)

type ExecutionWriteModel struct {
	eventstore.WriteModel

	Name             string
	ExecutionType    domain.ExecutionType
	URL              string
	Timeout          time.Duration
	Async            bool
	InterruptOnError bool

	State domain.ExecutionState
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
		case *execution.AddedEvent:
			wm.Name = e.Name
			wm.ExecutionType = e.ExecutionType
			wm.URL = e.URL
			wm.Timeout = e.Timeout
			wm.Async = e.Async
			wm.State = domain.ExecutionActive
		case *execution.ChangedEvent:
			if e.ExecutionType != nil {
				wm.ExecutionType = *e.ExecutionType
			}
			if e.URL != nil {
				wm.URL = *e.URL
			}
			if e.Timeout != nil {
				wm.Timeout = *e.Timeout
			}
			if e.Async != nil {
				wm.Async = *e.Async
			}
			if e.InterruptOnError != nil {
				wm.InterruptOnError = *e.InterruptOnError
			}
		case *action.RemovedEvent:
			wm.State = domain.ExecutionRemoved
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
		EventTypes(execution.AddedEventType,
			execution.ChangedEventType,
			execution.RemovedEventType).
		Builder()
}

func (wm *ExecutionWriteModel) NewChangedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	executionType domain.ExecutionType,
	url string,
	timeout time.Duration,
	async bool,
	interruptOnError bool,
) *execution.ChangedEvent {
	changes := make([]execution.Changes, 0)
	if wm.ExecutionType != executionType {
		changes = append(changes, execution.ChangeExecutionType(executionType))
	}
	if wm.URL != url {
		changes = append(changes, execution.ChangeURL(url))
	}
	if wm.Timeout != timeout {
		changes = append(changes, execution.ChangeTimeout(timeout))
	}
	if wm.Async != async {
		changes = append(changes, execution.ChangeAsync(async))
	}
	if wm.InterruptOnError != interruptOnError {
		changes = append(changes, execution.ChangeInterruptOnError(interruptOnError))
	}
	if len(changes) == 0 {
		return nil
	}
	return execution.NewChangedEvent(ctx, agg, changes)
}

func ExecutionAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, execution.AggregateType, execution.AggregateVersion)
}

func NewExecutionAggregate(id, resourceOwner string) *eventstore.Aggregate {
	return ExecutionAggregateFromWriteModel(&eventstore.WriteModel{
		AggregateID:   id,
		ResourceOwner: resourceOwner,
	})
}
