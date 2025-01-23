package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserGroupMemberWrite struct {
	eventstore.WriteModel

	GroupID string
	State   domain.GroupState
}

type UserGroupMemberWriteModel struct {
	UserGroupMemberWrite
}

func (wm *UserGroupMemberWrite) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.UserGroupAddedEvent:
			wm.GroupID = e.GroupID
			wm.State = domain.GroupStateActive
		case *user.UserGroupRemovedEvent:
			wm.GroupID = e.GroupID
			wm.State = domain.GroupStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func NewUserGroupMemberWriteModel(groupID, userID, resourceOwner string) *UserGroupMemberWriteModel {
	return &UserGroupMemberWriteModel{
		UserGroupMemberWrite{
			WriteModel: eventstore.WriteModel{
				AggregateID:   userID,
				ResourceOwner: resourceOwner,
			},
			GroupID: groupID,
		},
	}
}

func (wm *UserGroupMemberWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.UserGroupAddedEvent:
			if e.GroupID != wm.UserGroupMemberWrite.GroupID {
				continue
			}
			wm.UserGroupMemberWrite.AppendEvents(e)
		case *user.UserGroupRemovedEvent:
			if e.GroupID != wm.UserGroupMemberWrite.GroupID {
				continue
			}
			wm.UserGroupMemberWrite.AppendEvents(e)
		}
	}
}

func (wm *UserGroupMemberWriteModel) Reduce() error {
	return wm.UserGroupMemberWrite.Reduce()
}

func (wm *UserGroupMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.UserGroupAddedType,
			user.UserGroupRemovedType,
			user.UserGroupCascadeRemovedType).
		Builder()
}
