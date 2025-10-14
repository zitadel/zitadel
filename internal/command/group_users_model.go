package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type GroupUsersWriteModel struct {
	eventstore.WriteModel

	UserIDs         []string
	existingUserIDs []string
}

func NewGroupUsersWriteModel(resourceOwner, groupID string, userIDs []string) *GroupUsersWriteModel {
	return &GroupUsersWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
		UserIDs: userIDs,
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
				if !slices.Contains(g.existingUserIDs, userID) {
					g.existingUserIDs = append(g.existingUserIDs, userID)
				}
			}
		case *group.GroupUsersRemovedEvent:
			for _, userID := range e.UserIDs {
				i := slices.Index(g.existingUserIDs, userID)
				if i >= 0 {
					g.existingUserIDs = slices.Delete(g.existingUserIDs, i, i+1)
				}
			}
		}
	}
	return g.WriteModel.Reduce()
}

func (g *GroupUsersWriteModel) userIDsToAdd() []string {
	userIDsToAdd := make([]string, 0)
	for _, userID := range g.UserIDs {
		if !slices.Contains(g.existingUserIDs, userID) {
			userIDsToAdd = append(userIDsToAdd, userID)
		}
	}
	return userIDsToAdd
}
