package handler

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"time"
)

type ProjectRole struct {
	handler
	projectEvents *proj_event.ProjectEventstore
}

const (
	projectRoleTable = "management.project_roles"
)

func (p *ProjectRole) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *ProjectRole) ViewModel() string {
	return projectRoleTable
}

func (p *ProjectRole) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.ProjectQuery(sequence), nil
}

func (p *ProjectRole) Process(event *models.Event) (err error) {
	role := new(view_model.ProjectRoleView)
	switch event.Type {
	case es_model.ProjectRoleAdded:
		role.AppendEvent(event)
	case es_model.ProjectRoleChanged:
		err := role.SetData(event)
		if err != nil {
			return err
		}
		role, err = p.view.ProjectRoleByIDs(event.AggregateID, event.ResourceOwner, role.Key)
		if err != nil {
			return err
		}
		role.AppendEvent(event)
	case es_model.ProjectRoleRemoved:
		err := role.SetData(event)
		if err != nil {
			return err
		}
		err = p.removeRoleFromAllResourceowners(event, role)
	case es_model.ProjectGrantAdded:
		return p.addGrantRoles(event)
	case es_model.ProjectGrantChanged:
		err = p.removeRolesFromResourceowner(event)
		if err != nil {
			return err
		}
		return p.addGrantRoles(event)
	case es_model.ProjectGrantRemoved:
		return p.removeRolesFromResourceowner(event)
	default:
		return p.view.ProcessedProjectRoleSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectRole(role)
}

func (p *ProjectRole) removeRoleFromAllResourceowners(event *models.Event, role *view_model.ProjectRoleView) error {
	roles, err := p.view.ResourceOwnerProjectRolesByKey(event.AggregateID, event.ResourceOwner, role.Key)
	if err != nil {
		logging.LogWithFields("HANDL-slo03", "aggregateID", event.AggregateID, "ResourceOwner", event.ResourceOwner, "Key", role.Key).WithError(err).Warn("could not read roles to remove")
		return err
	}
	for _, r := range roles {
		err = p.view.DeleteProjectRole(r.ProjectID, r.OrgID, r.Key, event.Sequence)
		if err != nil {
			logging.LogWithFields("HANDL-kloa2", "aggregateID", event.AggregateID, "ResourceOwner", event.ResourceOwner, "OrgID", r.OrgID, "Key", role.Key).WithError(err).Warn("could not remove role")
			return err
		}
	}
	return nil
}

func (p *ProjectRole) removeRolesFromResourceowner(event *models.Event) error {
	roles, err := p.view.ResourceOwnerProjectRoles(event.AggregateID, event.ResourceOwner)
	if err != nil {
		logging.LogWithFields("HANDL-slo03", "aggregateID", event.AggregateID, "ResourceOwner", event.ResourceOwner, "Key").WithError(err).Warn("could not read roles to remove")
		return err
	}
	for _, r := range roles {
		err = p.view.DeleteProjectRole(r.ProjectID, r.OrgID, r.Key, event.Sequence)
		if err != nil {
			logging.LogWithFields("HANDL-kloa2", "aggregateID", event.AggregateID, "ResourceOwner", event.ResourceOwner, "OrgID", r.OrgID).WithError(err).Warn("could not remove role")
			return err
		}
	}
	return nil
}

func (p *ProjectRole) addGrantRoles(event *models.Event) error {
	project, err := p.projectEvents.ProjectByID(context.Background(), event.AggregateID)
	if err != nil {
		return err
	}

	grant := new(view_model.ProjectGrant)
	err = grant.SetData(event)
	if err != nil {
		return err
	}
	for _, roleKey := range grant.RoleKeys {
		role := getRoleFromProject(roleKey, project)
		projectRole := &view_model.ProjectRoleView{
			OrgID:         grant.GrantedOrgID,
			ProjectID:     event.AggregateID,
			Key:           roleKey,
			DisplayName:   role.DisplayName,
			Group:         role.Group,
			ResourceOwner: event.ResourceOwner,
			CreationDate:  event.CreationDate,
			Sequence:      event.Sequence,
		}
		err := p.view.PutProjectRole(projectRole)
		logging.LogWithFields("HANDL-sj3TG", "eventID", event.ID).OnError(err).Warn("could not save project role")
	}
	return nil
}

func getRoleFromProject(roleKey string, project *proj_model.Project) *proj_model.ProjectRole {
	for _, role := range project.Roles {
		if roleKey == role.Key {
			return role
		}
	}
	return nil
}

func (p *ProjectRole) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-lso9w", "id", event.AggregateID).WithError(err).Warn("something went wrong in project role handler")
	return spooler.HandleError(event, p.view.GetLatestProjectRoleFailedEvent, p.view.ProcessedProjectRoleFailedEvent, p.view.ProcessedProjectRoleSequence, p.errorCountUntilSkip)
}
