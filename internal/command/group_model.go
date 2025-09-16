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
	ID          string
	Description string

	State domain.GroupState
}

// Query constructs a search query for retrieving group-related events based on the GroupWriteModel attributes.
func (g GroupWriteModel) Query() *eventstore.SearchQueryBuilder {
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

// GroupAggregateFromWriteModel maps a WriteModel to a group-specific Aggregate using its type and version.
func GroupAggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, group.AggregateType, group.AggregateVersion)
}
