package command

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/target"
)

type TargetWriteModel struct {
	eventstore.WriteModel

	Name             string
	TargetType       domain.TargetType
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
			wm.TargetType = e.TargetType
			wm.URL = e.URL
			wm.Timeout = e.Timeout
			wm.Async = e.Async
			wm.State = domain.TargetActive
		case *target.ChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.TargetType != nil {
				wm.TargetType = *e.TargetType
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
		case *target.RemovedEvent:
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
	targetType *domain.TargetType,
	url *string,
	timeout *time.Duration,
	async *bool,
	interruptOnError *bool,
) *target.ChangedEvent {
	changes := make([]target.Changes, 0)
	if name != nil && wm.Name != *name {
		changes = append(changes, target.ChangeName(wm.Name, *name))
	}
	if targetType != nil && wm.TargetType != *targetType {
		changes = append(changes, target.ChangeTargetType(*targetType))
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

type TargetsExistsWriteModel struct {
	eventstore.WriteModel
	ids         []string
	existingIDs []string
}

func (e *TargetsExistsWriteModel) AllExists() bool {
	return len(e.ids) == len(e.existingIDs)
}

func NewTargetsExistsWriteModel(ids []string, resourceOwner string) *TargetsExistsWriteModel {
	return &TargetsExistsWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
		ids: ids,
	}
}

func (wm *TargetsExistsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *target.AddedEvent:
			if !slices.Contains(wm.existingIDs, e.Aggregate().ID) {
				wm.existingIDs = append(wm.existingIDs, e.Aggregate().ID)
			}
		case *target.RemovedEvent:
			i := slices.Index(wm.existingIDs, e.Aggregate().ID)
			if i >= 0 {
				wm.existingIDs = slices.Delete(wm.existingIDs, i, i+1)
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *TargetsExistsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(target.AggregateType).
		AggregateIDs(wm.ids...).
		EventTypes(target.AddedEventType,
			target.RemovedEventType).
		Builder()
}

func TargetAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          target.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       target.AggregateVersion,
	}
}
