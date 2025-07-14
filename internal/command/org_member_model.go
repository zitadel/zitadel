package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
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

func (wm *OrgMemberWriteModel) AppendEvents(events ...eventstore.Event) {
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
		case *org.MemberCascadeRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberCascadeRemovedEvent)
		case *org.OrgRemovedEvent:
			wm.MemberWriteModel.AppendEvents(e)
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
			org.MemberRemovedEventType,
			org.MemberCascadeRemovedEventType,
			org.OrgRemovedEventType).
		Builder()
}
