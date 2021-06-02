package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgMemberWriteModel struct {
	MemberWriteModel
}

func NewOrgMemberWriteModel(orgID, userID string) *OrgMemberWriteModel {
	return &OrgMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			UserID: userID,
		},
	}
}

func (wm *OrgMemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *org.MemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *org.MemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		}
	}
}

func (wm *OrgMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *OrgMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		EventTypes(
			org.MemberAddedEventType,
			org.MemberChangedEventType,
			org.MemberRemovedEventType).
		Builder()
}
