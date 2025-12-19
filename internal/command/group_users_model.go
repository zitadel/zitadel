package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	groupusers "github.com/zitadel/zitadel/internal/repository/group_users"
)

type GroupUserWriteModel struct {
	eventstore.WriteModel
	UserID     string
	State      domain.GroupUserState
	Attributes []string
}

func NewGroupUserWriteModel(groupID, userID string, resourceOwner string, attributes []string) *GroupUserWriteModel {
	return &GroupUserWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   groupID,
			ResourceOwner: resourceOwner,
		},
		UserID:     userID,
		Attributes: attributes,
	}
}

func (wm *GroupUserWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *groupusers.GroupUserAddedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.Attributes = e.Attributes
			wm.WriteModel.AppendEvents(e)
		case *groupusers.GroupUserChangedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.Attributes = e.Attributes
			wm.WriteModel.AppendEvents(e)
		case *groupusers.GroupUserRemovedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *groupusers.GroupUserCascadeRemovedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *GroupUserWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *GroupUserWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(groupusers.AddedEventType,
			groupusers.ChangedEventType,
			groupusers.RemovedEventType,
			groupusers.CascadeRemovedEventType).
		Builder()
}
