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

	UserIDs         []string
	existingUserIDs map[string]struct{}
}

// NewGroupWriteModel initializes a new instance of GroupWriteModel from the given Group.
func NewGroupWriteModel(groupID, orgID string, userIDs []string) *GroupWriteModel {
	return &GroupWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: orgID,
		},
		UserIDs:         userIDs,
		existingUserIDs: make(map[string]struct{}),
	}
}

func (g *GroupWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &g.WriteModel
}

// Query constructs a search query for retrieving group-related events based on the GroupWriteModel attributes.
func (g *GroupWriteModel) Query() *eventstore.SearchQueryBuilder {
	eventTypes := []eventstore.EventType{
		group.GroupAddedEventType,
		group.GroupChangedEventType,
		group.GroupRemovedEventType,
	}
	if g.UserIDs != nil {
		eventTypes = append(eventTypes, group.GroupUsersAddedEventType, group.GroupUsersRemovedEventType)
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(g.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(g.AggregateID).
		EventTypes(eventTypes...).Builder()
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
			if e.Name != nil {
				g.Name = *e.Name
			}
			if e.Description != nil {
				g.Description = *e.Description
			}
		case *group.GroupRemovedEvent:
			g.State = domain.GroupStateRemoved
		case *group.GroupUsersAddedEvent:
			for _, userID := range e.UserIDs {
				g.existingUserIDs[userID] = struct{}{}
			}
		case *group.GroupUsersRemovedEvent:
			for _, userID := range e.UserIDs {
				delete(g.existingUserIDs, userID)
			}
		}
	}
	return g.WriteModel.Reduce()
}

func (g *GroupWriteModel) NewChangedEvent(ctx context.Context, agg *eventstore.Aggregate, name, description *string) *group.GroupChangedEvent {
	changes := make([]group.GroupChanges, 0)
	oldName := ""

	if name != nil && g.Name != *name {
		oldName = g.Name
		changes = append(changes, group.ChangeName(name))
	}
	if description != nil && g.Description != *description {
		changes = append(changes, group.ChangeDescription(description))
	}
	if len(changes) == 0 {
		return nil
	}

	return group.NewGroupChangedEvent(ctx, agg, oldName, changes)
}

// GroupAggregateFromWriteModel maps a WriteModel to a group-specific Aggregate using its type and version.
func GroupAggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, group.AggregateType, group.AggregateVersion)
}
