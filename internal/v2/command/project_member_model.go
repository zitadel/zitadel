package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ProjectMemberWriteModel struct {
	MemberWriteModel
}

func NewProjectMemberWriteModel(projectID, userID, resourceOwner string) *ProjectMemberWriteModel {
	return &ProjectMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			UserID: userID,
		},
	}
}

func (wm *ProjectMemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ProjectMemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *project.ProjectMemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *project.ProjectMemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		}
	}
}

func (wm *ProjectMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *ProjectMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
