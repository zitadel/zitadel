package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceMemberWriteModel struct {
	MemberWriteModel
}

func NewInstanceMemberWriteModel(instanceID, userID string) *InstanceMemberWriteModel {
	return &InstanceMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			UserID: userID,
		},
	}
}

func (wm *InstanceMemberWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.MemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *instance.MemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *instance.MemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		case *instance.MemberCascadeRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberCascadeRemovedEvent)
		}
	}
}

func (wm *InstanceMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *InstanceMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		EventTypes(
			instance.MemberAddedEventType,
			instance.MemberChangedEventType,
			instance.MemberRemovedEventType,
			instance.MemberCascadeRemovedEventType).
		Builder()
}
