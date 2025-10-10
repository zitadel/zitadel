package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type GroupUserWriteModel struct {
	eventstore.WriteModel

	UserID string
	State  domain.GroupUserState
}

func (g *GroupUserWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &g.WriteModel
}

func NewGroupUserWriteModel(resourceOwner, groupID, userID string) *GroupUserWriteModel {
	return &GroupUserWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
		UserID: userID,
	}
}

func (g *GroupUserWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(g.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(g.AggregateID).
		EventTypes(group.GroupUserAddedEventType,
			group.GroupUserRemovedEventType).Builder()
}

func (g *GroupUserWriteModel) Reduce() error {
	for _, event := range g.Events {
		switch e := event.(type) {
		case *group.GroupUserAddedEvent:
			g.AggregateID = e.ID
			g.UserID = e.UserID
			g.State = domain.GroupUserStateActive
		case *group.GroupUserRemovedEvent:
			g.State = domain.GroupUserStateRemoved
		}
	}
	return g.WriteModel.Reduce()
}
