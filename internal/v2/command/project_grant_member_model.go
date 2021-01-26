package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ProjectGrantMemberWriteModel struct {
	eventstore.WriteModel

	GrantID string
	UserID  string
	Roles   []string

	State domain.MemberState
}

func NewProjectGrantMemberWriteModel(projectID, userID, grantID, resourceOwner string) *ProjectGrantMemberWriteModel {
	return &ProjectGrantMemberWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		UserID: userID,
	}
}

func (wm *ProjectGrantMemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ProjectGrantMemberAddedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectGrantMemberChangedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectGrantMemberRemovedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.GrantRemovedEvent:
			if e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ProjectGrantMemberWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ProjectGrantMemberAddedEvent:
			wm.Roles = e.Roles
			wm.State = domain.MemberStateActive
		case *project.ProjectGrantMemberChangedEvent:
			wm.Roles = e.Roles
		case *project.ProjectGrantMemberRemovedEvent:
			wm.State = domain.MemberStateRemoved
		case *project.GrantRemovedEvent, *project.ProjectRemovedEvent:
			wm.State = domain.MemberStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ProjectGrantMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			project.ProjectGrantMemberAddedType,
			project.ProjectGrantMemberChangedType,
			project.ProjectGrantMemberRemovedType,
			project.GrantRemovedType,
			project.ProjectRemovedType)
}
