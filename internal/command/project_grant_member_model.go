package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
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
		UserID:  userID,
		GrantID: grantID,
	}
}

func (wm *ProjectGrantMemberWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.GrantMemberAddedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.GrantMemberChangedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.GrantMemberRemovedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.GrantMemberCascadeRemovedEvent:
			if e.UserID != wm.UserID || e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.GrantRemovedEvent:
			if e.GrantID != wm.GrantID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ProjectGrantMemberWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.GrantMemberAddedEvent:
			wm.Roles = e.Roles
			wm.State = domain.MemberStateActive
			wm.ResourceOwner = e.Aggregate().ResourceOwner
		case *project.GrantMemberChangedEvent:
			wm.Roles = e.Roles
		case *project.GrantMemberRemovedEvent:
			wm.State = domain.MemberStateRemoved
		case *project.GrantMemberCascadeRemovedEvent:
			wm.State = domain.MemberStateRemoved
		case *project.GrantRemovedEvent, *project.ProjectRemovedEvent:
			wm.State = domain.MemberStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ProjectGrantMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.GrantMemberAddedType,
			project.GrantMemberChangedType,
			project.GrantMemberRemovedType,
			project.GrantMemberCascadeRemovedType,
			project.GrantRemovedType,
			project.ProjectRemovedType).
		Builder()
	return query
}
