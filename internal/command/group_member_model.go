package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type GroupMemberWriteModel struct {
	eventstore.WriteModel

	UserID string
	State  domain.GroupMemberState
}

type GroupMemberWrite struct {
	GroupMemberWriteModel
}

func NewGroupMemberWriteModel(groupID, userID, resourceOwner string) *GroupMemberWrite {
	return &GroupMemberWrite{
		GroupMemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   groupID,
				ResourceOwner: resourceOwner,
			},
			UserID: userID,
		},
	}
}

func (wm *GroupMemberWrite) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *group.GroupMemberAddedEvent:
			if e.UserID != wm.GroupMemberWriteModel.UserID {
				continue
			}
			wm.GroupMemberWriteModel.AppendEvents(&e.GroupMemberAddedEvent)
		case *group.GroupMemberChangedEvent:
			if e.UserID != wm.GroupMemberWriteModel.UserID {
				continue
			}
			wm.GroupMemberWriteModel.AppendEvents(&e.GroupMemberChangedEvent)
		case *group.GroupMemberRemovedEvent:
			if e.UserID != wm.GroupMemberWriteModel.UserID {
				continue
			}
			wm.GroupMemberWriteModel.AppendEvents(&e.GroupMemberRemovedEvent)
		case *group.GroupMemberCascadeRemovedEvent:
			if e.UserID != wm.GroupMemberWriteModel.UserID {
				continue
			}
			wm.GroupMemberWriteModel.AppendEvents(&e.GroupMemberCascadeRemovedEvent)
		}
	}
}

func (wm *GroupMemberWrite) Reduce() error {
	return wm.GroupMemberWriteModel.Reduce()
}

func (wm *GroupMemberWrite) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(wm.GroupMemberWriteModel.AggregateID).
		EventTypes(group.MemberAddedType,
			group.MemberChangedType,
			group.MemberRemovedType,
			group.MemberCascadeRemovedType).
		Builder()
}
