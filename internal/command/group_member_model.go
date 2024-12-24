package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

type GroupMemberWriteModel struct {
	MemberWriteModel
}

func NewGroupMemberWriteModel(groupID, userID, resourceOwner string) *GroupMemberWriteModel {
	return &GroupMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   groupID,
				ResourceOwner: resourceOwner,
			},
			UserID: userID,
		},
	}
}

func (wm *GroupMemberWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *group.MemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *group.MemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *group.MemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		case *group.MemberCascadeRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberCascadeRemovedEvent)
		}
	}
}

func (wm *GroupMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *GroupMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		EventTypes(group.MemberAddedType,
			group.MemberChangedType,
			group.MemberRemovedType,
			group.MemberCascadeRemovedType).
		Builder()
}
