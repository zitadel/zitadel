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

type Application struct {
	handler
	projectEvents *proj_event.ProjectEventstore
}

const (
	applicationTable = "management.applications"
)

func (p *Application) ViewModel() string {
	return applicationTable
}

func (p *Application) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestApplicationSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *Application) Reduce(event *models.Event) (err error) {
	app := new(view_model.ApplicationView)
	switch event.Type {
	case es_model.ApplicationAdded:
		app.AppendEvent(event)
	case es_model.ApplicationChanged,
		es_model.OIDCConfigAdded,
		es_model.OIDCConfigChanged,
		es_model.ApplicationDeactivated,
		es_model.ApplicationReactivated:
		err := app.SetData(event)
		if err != nil {
			return err
		}
		app, err = p.view.ApplicationByID(app.ID)
		if err != nil {
			return err
		}
		app.AppendEvent(event)
	case es_model.ApplicationRemoved:
		err := app.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteApplication(app.ID, event.Sequence)
	case es_model.ProjectRemoved:
		return p.view.DeleteApplicationsByProjectID(event.AggregateID)
	default:
		return p.view.ProcessedApplicationSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutApplication(app)
}

func (p *Application) OnError(event *models.Event, spoolerError error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerError).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerError, p.view.GetLatestApplicationFailedEvent, p.view.ProcessedApplicationFailedEvent, p.view.ProcessedApplicationSequence, p.errorCountUntilSkip)
}
