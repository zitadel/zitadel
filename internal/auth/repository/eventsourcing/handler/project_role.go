package handler

import (
	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_events "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
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

func (_ *ProjectRole) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.ProjectAggregate}
}

func (p *ProjectRole) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectRole) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return nil, err
	}
	return proj_events.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *ProjectRole) Reduce(event *es_models.Event) (err error) {
	role := new(view_model.ProjectRoleView)
	switch event.Type {
	case model.ProjectRoleAdded:
		err = role.AppendEvent(event)
	case model.ProjectRoleChanged:
		err = role.SetData(event)
		if err != nil {
			return err
		}
		role, err = p.view.ProjectRoleByIDs(event.AggregateID, event.ResourceOwner, role.Key)
		if err != nil {
			return err
		}
		err = role.AppendEvent(event)
	case model.ProjectRoleRemoved:
		err = role.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectRole(event.AggregateID, event.ResourceOwner, role.Key, event.Sequence, event.CreationDate)
	case model.ProjectRemoved:
		return p.view.DeleteProjectRolesByProjectID(event.AggregateID)
	default:
		return p.view.ProcessedProjectRoleSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectRole(role, event.CreationDate)
}

func (p *ProjectRole) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-lso9w", "id", event.AggregateID).WithError(err).Warn("something went wrong in project role handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectRoleFailedEvent, p.view.ProcessedProjectRoleFailedEvent, p.view.ProcessedProjectRoleSequence, p.errorCountUntilSkip)
}

func (p *ProjectRole) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectRoleSpoolerRunTimestamp)
}
