package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type GroupUsersWriteModel struct {
	eventstore.WriteModel

	UserIDs         []string
	existingUserIDs map[string]struct{}
}

func NewGroupUsersWriteModel(resourceOwner, groupID string, userIDs []string) *GroupUsersWriteModel {
	return &GroupUsersWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
		UserIDs:         userIDs,
		existingUserIDs: make(map[string]struct{}),
	}
}

func (g *GroupUsersWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(g.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(g.AggregateID).
		EventTypes(group.GroupUsersAddedEventType,
			group.GroupUsersRemovedEventType).Builder()
}

func (g *GroupUsersWriteModel) Reduce() error {
	for _, event := range g.Events {
		switch e := event.(type) {
		case *group.GroupUsersAddedEvent:
			for _, userID := range e.UserIDs {
				if _, ok := g.existingUserIDs[userID]; !ok {
					g.existingUserIDs[userID] = struct{}{}
				}

			}
		case *group.GroupUsersRemovedEvent:
			for _, userID := range e.UserIDs {
				if _, ok := g.existingUserIDs[userID]; ok {
					delete(g.existingUserIDs, userID)
				}
			}
		}
	}
	return g.WriteModel.Reduce()
}

func (g *GroupUsersWriteModel) userIDsToAdd() []string {
	userIDsToAdd := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; !ok {
			userIDsToAdd = append(userIDsToAdd, userID)
		}
	}
	return userIDsToAdd
}

func (g *GroupUsersWriteModel) userIDsToRemove() []string {
	userIDsToRemove := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; ok {
			userIDsToRemove = append(userIDsToRemove, userID)
		}
	}
	return userIDsToRemove
}
