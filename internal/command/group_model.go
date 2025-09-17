package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

// GroupWriteModel represents the write-model for a group.
type GroupWriteModel struct {
	eventstore.WriteModel

	Name        string
	Description string

	State domain.GroupState
}

// NewGroupWriteModel initializes a new instance of GroupWriteModel from the given Group.
func NewGroupWriteModel(group *domain.Group) *GroupWriteModel {
	return &GroupWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   group.AggregateID,
			ResourceOwner: group.ResourceOwner,
		},
		Name:        group.Name,
		Description: group.Description,
	}
}

// Query constructs a search query for retrieving group-related events based on the GroupWriteModel attributes.
func (g *GroupWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(g.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(g.AggregateID).
		EventTypes(
			group.GroupAddedEventType,
			group.GroupChangedEventType,
			group.GroupRemovedEventType).Builder()
}

func (g *GroupWriteModel) Reduce() error {
	for _, event := range g.Events {
		switch e := event.(type) {
		case *group.GroupAddedEvent:
			g.AggregateID = e.ID
			g.Name = e.Name
			g.Description = e.Description
			g.State = domain.GroupStateActive
		case *group.GroupChangedEvent:
			g.Name = e.Name
			g.Description = e.Description
		case *group.GroupRemovedEvent:
			g.State = domain.GroupStateRemoved
		}
	}
	return g.WriteModel.Reduce()
}

// GroupAggregateFromWriteModel maps a WriteModel to a group-specific Aggregate using its type and version.
func GroupAggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, group.AggregateType, group.AggregateVersion)
}
