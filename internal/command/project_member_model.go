package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
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

func (wm *ProjectMemberWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.MemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *project.MemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *project.MemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		case *project.MemberCascadeRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberCascadeRemovedEvent)
		}
	}
}

func (wm *ProjectMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *ProjectMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		EventTypes(project.MemberAddedType,
			project.MemberChangedType,
			project.MemberRemovedType,
			project.MemberCascadeRemovedType).
		Builder()
}
