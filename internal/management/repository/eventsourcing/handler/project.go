package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type Project struct {
	handler
	eventstore eventstore.Eventstore
}

const (
	projectTable = "management.projects"
)

func (p *Project) ViewModel() string {
	return projectTable
}

func (p *Project) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectSequence()
	if err != nil {
		return nil, err
	}
	return proj_event.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *Project) Reduce(event *models.Event) (err error) {
	project := new(view_model.ProjectView)
	switch event.Type {
	case es_model.ProjectAdded:
		err = project.AppendEvent(event)
	case es_model.ProjectChanged,
		es_model.ProjectDeactivated,
		es_model.ProjectReactivated:
		project, err = p.view.ProjectByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = project.AppendEvent(event)
	case es_model.ProjectRemoved:
		return p.view.DeleteProject(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedProjectSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutProject(project, event.CreationDate)
}

func (p *Project) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-dLsop3", "id", event.AggregateID).WithError(err).Warn("something went wrong in projecthandler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectFailedEvent, p.view.ProcessedProjectFailedEvent, p.view.ProcessedProjectSequence, p.errorCountUntilSkip)
}

func (p *Project) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectSpoolerRunTimestamp)
}
