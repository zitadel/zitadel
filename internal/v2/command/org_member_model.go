package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type OrgMemberWriteModel struct {
	MemberWriteModel
}

func NewOrgMemberWriteModel(orgID, userID string) *OrgMemberWriteModel {
	return &OrgMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: orgID,
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID)
}
