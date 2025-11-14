package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ProjectGrantWriteModel struct {
	eventstore.WriteModel

	GrantID      string
	GrantedOrgID string
	RoleKeys     []string
	State        domain.ProjectGrantState

	FoundGrantID string
}

func NewProjectGrantWriteModel(grantID, grantedOrgID, projectID, resourceOwner string) *ProjectGrantWriteModel {
	return &ProjectGrantWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		// Always either the grantID or the grantedOrgID is provided
		GrantID:      grantID,
		GrantedOrgID: grantedOrgID,
	}
}

func (wm *ProjectGrantWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.GrantAddedEvent:
			if projectGrantExists(wm.GrantID, wm.GrantedOrgID, e.GrantID, e.GrantedOrgID) {
				wm.FoundGrantID = e.GrantID
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantChangedEvent:
			if projectGrantEqual(wm.FoundGrantID, e.GrantID) {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantCascadeChangedEvent:
			if projectGrantEqual(wm.FoundGrantID, e.GrantID) {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantDeactivateEvent:
			if projectGrantEqual(wm.FoundGrantID, e.GrantID) {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantReactivatedEvent:
			if projectGrantEqual(wm.FoundGrantID, e.GrantID) {
				wm.WriteModel.AppendEvents(e)
			}
		case *project.GrantRemovedEvent:
			if projectGrantEqual(wm.FoundGrantID, e.GrantID) {
				wm.FoundGrantID = ""
				wm.WriteModel.AppendEvents(e)
			}
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func projectGrantExists(requiredGrantID, requiredGrantedOrgID, grantID, grantedOrgID string) bool {
	// either grantID or grantedOrgID is provided and equal
	return projectGrantEqual(requiredGrantID, grantID) ||
		(requiredGrantedOrgID != "" && grantedOrgID == requiredGrantedOrgID)
}

func projectGrantEqual(requiredGrantID, grantID string) bool {
	// grantID is provided and equal
	return requiredGrantID != "" && grantID == requiredGrantID
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.GrantAddedType,
			project.GrantChangedType,
			project.GrantCascadeChangedType,
			project.GrantDeactivatedType,
			project.GrantReactivatedType,
			project.GrantRemovedType,
			project.ProjectRemovedType).
		Builder()
}

type ProjectGrantPreConditionReadModel struct {
	eventstore.WriteModel

	ProjectResourceOwner string
	ProjectID            string
	GrantedOrgID         string
	ProjectExists        bool
	GrantedOrgExists     bool
	ExistingRoleKeys     []string
}

func NewProjectGrantPreConditionReadModel(projectID, grantedOrgID, resourceOwner string) *ProjectGrantPreConditionReadModel {
	return &ProjectGrantPreConditionReadModel{
		WriteModel:           eventstore.WriteModel{},
		ProjectResourceOwner: resourceOwner,
		ProjectID:            projectID,
		GrantedOrgID:         grantedOrgID,
	}
}

func (wm *ProjectGrantPreConditionReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ProjectAddedEvent:
			if wm.ProjectResourceOwner == "" {
				wm.ProjectResourceOwner = e.Aggregate().ResourceOwner
			}
			if wm.ProjectResourceOwner != e.Aggregate().ResourceOwner {
				continue
			}
			wm.ProjectExists = true
		case *project.ProjectRemovedEvent:
			if wm.ProjectResourceOwner != e.Aggregate().ResourceOwner {
				continue
			}
			wm.ProjectResourceOwner = ""
			wm.ProjectExists = false
		case *project.RoleAddedEvent:
			if e.Aggregate().ResourceOwner != wm.ProjectResourceOwner {
				continue
			}
			wm.ExistingRoleKeys = append(wm.ExistingRoleKeys, e.Key)
		case *project.RoleRemovedEvent:
			if e.Aggregate().ResourceOwner != wm.ProjectResourceOwner {
				continue
			}
			wm.ExistingRoleKeys = slices.DeleteFunc(wm.ExistingRoleKeys, func(key string) bool {
				return key == e.Key
			})
		case *org.OrgAddedEvent:
			wm.GrantedOrgExists = true
		case *org.OrgRemovedEvent:
			wm.GrantedOrgExists = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ProjectGrantPreConditionReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.ProjectID).
		EventTypes(
			project.ProjectAddedType,
			project.ProjectRemovedType,
			project.RoleAddedType,
			project.RoleRemovedType).
		Or().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.GrantedOrgID).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgRemovedEventType).
		Builder()

	return query
}
