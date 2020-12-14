package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type Application struct {
	handler
}

const (
	applicationTable = "authz.applications"
)

func (a *Application) ViewModel() string {
	return applicationTable
}

func (a *Application) EventQuery() (*models.SearchQuery, error) {
	sequence, err := a.view.GetLatestApplicationSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.ProjectQuery(sequence.CurrentSequence), nil
}

func (a *Application) Reduce(event *models.Event) (err error) {
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
		app, err = a.view.ApplicationByID(event.AggregateID, app.ID)
		if err != nil {
			return err
		}
		app.AppendEvent(event)
	case es_model.ApplicationRemoved:
		err := app.SetData(event)
		if err != nil {
			return err
		}
		return a.view.DeleteApplication(app.ID, event.Sequence, event.CreationDate)
	default:
		return a.view.ProcessedApplicationSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return a.view.PutApplication(app, event.CreationDate)
}

func (a *Application) OnError(event *models.Event, spoolerError error) error {
	logging.LogWithFields("SPOOL-sjZw", "id", event.AggregateID).WithError(spoolerError).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerError, a.view.GetLatestApplicationFailedEvent, a.view.ProcessedApplicationFailedEvent, a.view.ProcessedApplicationSequence, a.errorCountUntilSkip)
}

func (a *Application) OnSuccess() error {
	return spooler.HandleSuccess(a.view.UpdateApplicationSpoolerRunTimestamp)
}
