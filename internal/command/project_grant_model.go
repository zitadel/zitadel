package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectGrantWriteModel struct {
	eventstore.WriteModel

	GrantID      string
	GrantedOrgID string
	RoleKeys     []string
	State        domain.ProjectGrantState
}

func NewProjectGrantWriteModel(grantID, projectID, resourceOwner string) *ProjectGrantWriteModel {
	return &ProjectGrantWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		GrantID: grantID,
	}
}

func (wm *ProjectGrantWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.GrantAddedEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantChangedEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantCascadeChangedEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantDeactivateEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantReactivatedEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantRemovedEvent:
			if e.GrantID == wm.GrantID {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ProjectGrantWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.GrantAddedEvent:
			wm.GrantID = e.GrantID
			wm.GrantedOrgID = e.GrantedOrgID
			wm.RoleKeys = e.RoleKeys
			wm.State = domain.ProjectGrantStateActive
		case *project.GrantChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *project.GrantCascadeChangedEvent:
			wm.RoleKeys = e.RoleKeys
		case *project.GrantDeactivateEvent:
			if wm.State == domain.ProjectGrantStateRemoved {
				continue
			}
			wm.State = domain.ProjectGrantStateInactive
		case *project.GrantReactivatedEvent:
			if wm.State == domain.ProjectGrantStateRemoved {
				continue
			}
			wm.State = domain.ProjectGrantStateActive
		case *project.GrantRemovedEvent:
			wm.State = domain.ProjectGrantStateRemoved
		case *project.ProjectRemovedEvent:
			wm.State = domain.ProjectGrantStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ProjectGrantWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.GrantAddedType,
			project.GrantChangedType,
			project.GrantCascadeChangedType,
			project.GrantDeactivatedType,
			project.GrantReactivatedType,
			project.GrantRemovedType,
			project.ProjectRemovedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
