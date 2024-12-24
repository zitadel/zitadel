package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type GroupWriteModel struct {
	eventstore.WriteModel

	Name        string
	Description string
	State       domain.GroupState
}

func NewGroupWriteModel(groupID string, resourceOwner string) *GroupWriteModel {
	return &GroupWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *GroupWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *group.GroupAddedEvent:
			wm.Name = e.Name
			wm.Description = e.Description
			wm.State = domain.GroupStateActive
		case *group.GroupChangeEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Description != nil {
				wm.Description = *e.Description
			}
		case *group.GroupDeactivatedEvent:
			if wm.State == domain.GroupStateRemoved {
				continue
			}
			wm.State = domain.GroupStateInactive
		case *project.ProjectReactivatedEvent:
			if wm.State == domain.GroupStateRemoved {
				continue
			}
			wm.State = domain.GroupStateActive
		case *project.ProjectRemovedEvent:
			wm.State = domain.GroupStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GroupWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(group.GroupAddedType,
			group.GroupChangedType,
			group.GroupDeactivatedType,
			group.GroupReactivatedType,
			group.GroupRemovedType).
		Builder()
}

func (wm *GroupWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	description string,
) (*group.GroupChangeEvent, bool, error) {
	changes := make([]group.GroupChanges, 0)
	var err error

	oldName := ""
	oldDescription := ""
	if wm.Name != name {
		oldName = wm.Name
		changes = append(changes, group.ChangeName(name))
	}
	if wm.Description != description {
		oldDescription = wm.Description
		changes = append(changes, group.ChangeDescription(description))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := group.NewGroupChangeEvent(ctx, aggregate, oldName, oldDescription, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func GroupAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, group.AggregateType, group.AggregateVersion)
}
