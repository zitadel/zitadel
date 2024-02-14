package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/target"
)

type TargetWriteModel struct {
	eventstore.WriteModel

	Name             string
	ExecutionType    domain.TargetType
	URL              string
	Timeout          time.Duration
	Async            bool
	InterruptOnError bool

	State domain.TargetState
}

func NewTargetWriteModel(id string, resourceOwner string) *TargetWriteModel {
	return &TargetWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
	}
}

func (wm *TargetWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *target.AddedEvent:
			wm.Name = e.Name
			wm.ExecutionType = e.ExecutionType
			wm.URL = e.URL
			wm.Timeout = e.Timeout
			wm.Async = e.Async
			wm.State = domain.TargetActive
		case *target.ChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
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
			wm.State = domain.TargetRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *TargetWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(target.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(target.AddedEventType,
			target.ChangedEventType,
			target.RemovedEventType).
		Builder()
}

func (wm *TargetWriteModel) NewChangedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	name *string,
	executionType *domain.TargetType,
	url *string,
	timeout *time.Duration,
	async *bool,
	interruptOnError *bool,
) *target.ChangedEvent {
	changes := make([]target.Changes, 0)
	if name != nil && wm.Name != *name {
		changes = append(changes, target.ChangeName(wm.Name, *name))
	}
	if executionType != nil && wm.ExecutionType != *executionType {
		changes = append(changes, target.ChangeExecutionType(*executionType))
	}
	if url != nil && wm.URL != *url {
		changes = append(changes, target.ChangeURL(*url))
	}
	if timeout != nil && wm.Timeout != *timeout {
		changes = append(changes, target.ChangeTimeout(*timeout))
	}
	if async != nil && wm.Async != *async {
		changes = append(changes, target.ChangeAsync(*async))
	}
	if interruptOnError != nil && wm.InterruptOnError != *interruptOnError {
		changes = append(changes, target.ChangeInterruptOnError(*interruptOnError))
	}
	if len(changes) == 0 {
		return nil
	}
	return target.NewChangedEvent(ctx, agg, changes)
}

func TargetAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, target.AggregateType, target.AggregateVersion)
}

func NewExecutionAggregate(id, resourceOwner string) *eventstore.Aggregate {
	return TargetAggregateFromWriteModel(&eventstore.WriteModel{
		AggregateID:   id,
		ResourceOwner: resourceOwner,
	})
}
