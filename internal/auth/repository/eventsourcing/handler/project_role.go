package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type ProjectRole struct {
	handler
	projectEvents *proj_event.ProjectEventstore
}

const (
	projectRoleTable = "auth.project_roles"
)

func (p *ProjectRole) ViewModel() string {
	return projectRoleTable
}

func (p *ProjectRole) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *ProjectRole) Reduce(event *models.Event) (err error) {
	role := new(view_model.ProjectRoleView)
	switch event.Type {
	case es_model.ProjectRoleAdded:
		err = role.AppendEvent(event)
	case es_model.ProjectRoleChanged:
		err = role.SetData(event)
		if err != nil {
			return err
		}
		role, err = p.view.ProjectRoleByIDs(event.AggregateID, event.ResourceOwner, role.Key)
		if err != nil {
			return err
		}
		err = role.AppendEvent(event)
	case es_model.ProjectRoleRemoved:
		err = role.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectRole(event.AggregateID, event.ResourceOwner, role.Key, event.Sequence, event.CreationDate)
	case es_model.ProjectRemoved:
		return p.view.DeleteProjectRolesByProjectID(event.AggregateID)
	default:
		return p.view.ProcessedProjectRoleSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectRole(role, event.CreationDate)
}

func (p *ProjectRole) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-lso9w", "id", event.AggregateID).WithError(err).Warn("something went wrong in project role handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectRoleFailedEvent, p.view.ProcessedProjectRoleFailedEvent, p.view.ProcessedProjectRoleSequence, p.errorCountUntilSkip)
}

func (p *ProjectRole) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectRoleSpoolerRunTimestamp)
}
