package handler

import (
	"context"

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
	applicationTable = "auth.applications"
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

func (a *Application) Reduce(event *models.Event) (err error) {
	app := new(view_model.ApplicationView)
	switch event.Type {
	case es_model.ApplicationAdded:
		project, err := a.projectEvents.ProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		app.ProjectRoleCheck = project.ProjectRoleCheck
		app.ProjectRoleAssertion = project.ProjectRoleAssertion

		err = app.AppendEvent(event)
	case es_model.ApplicationChanged,
		es_model.OIDCConfigAdded,
		es_model.OIDCConfigChanged,
		es_model.ApplicationDeactivated,
		es_model.ApplicationReactivated:
		err = app.SetData(event)
		if err != nil {
			return err
		}
		app, err = a.view.ApplicationByID(event.AggregateID, app.ID)
		if err != nil {
			return err
		}
		err = app.AppendEvent(event)
	case es_model.ApplicationRemoved:
		err = app.SetData(event)
		if err != nil {
			return err
		}
		return a.view.DeleteApplication(app.ID, event.Sequence)
	case es_model.ProjectChanged:
		apps, err := a.view.ApplicationsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			return a.view.ProcessedApplicationSequence(event.Sequence)
		}
		for _, app := range apps {
			if err := app.AppendEvent(event); err != nil {
				return err
			}
		}
		return a.view.PutApplications(apps, event.Sequence)
	case es_model.ProjectRemoved:
		return a.view.DeleteApplicationsByProjectID(event.AggregateID)
	default:
		return a.view.ProcessedApplicationSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return a.view.PutApplication(app)
}

func (p *Application) OnError(event *models.Event, spoolerError error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerError).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerError, p.view.GetLatestApplicationFailedEvent, p.view.ProcessedApplicationFailedEvent, p.view.ProcessedApplicationSequence, p.errorCountUntilSkip)
}
