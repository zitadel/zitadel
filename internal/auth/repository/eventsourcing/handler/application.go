package handler

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v1"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	proj_model "github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

const (
	applicationTable = "auth.applications"
)

type Application struct {
	handler
	subscription *v1.Subscription
}

func newApplication(handler handler) *Application {
	h := &Application{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (a *Application) subscribe() {
	a.subscription = a.es.Subscribe(a.AggregateTypes()...)
	go func() {
		for event := range a.subscription.Events {
			query.ReduceEvent(a, event)
		}
	}()
}

func (a *Application) ViewModel() string {
	return applicationTable
}

func (_ *Application) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.ProjectAggregate}
}

func (a *Application) CurrentSequence() (uint64, error) {
	sequence, err := a.view.GetLatestApplicationSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (a *Application) EventQuery() (*models.SearchQuery, error) {
	sequence, err := a.view.GetLatestApplicationSequence()
	if err != nil {
		return nil, err
	}
	return proj_view.ProjectQuery(sequence.CurrentSequence), nil
}

func (a *Application) Reduce(event *models.Event) (err error) {
	app := new(view_model.ApplicationView)
	switch event.Type {
	case es_model.ApplicationAdded:
		project, err := a.getProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		app.ProjectRoleCheck = project.ProjectRoleCheck
		app.ProjectRoleAssertion = project.ProjectRoleAssertion

		err = app.AppendEvent(event)
	case es_model.ApplicationChanged,
		es_model.OIDCConfigAdded,
		es_model.OIDCConfigChanged,
		es_model.APIConfigAdded,
		es_model.APIConfigChanged,
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
		return a.view.DeleteApplication(app.ID, event)
	case es_model.ProjectChanged:
		apps, err := a.view.ApplicationsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			return a.view.ProcessedApplicationSequence(event)
		}
		for _, app := range apps {
			if err := app.AppendEvent(event); err != nil {
				return err
			}
		}
		return a.view.PutApplications(apps, event)
	case es_model.ProjectRemoved:
		err = a.view.DeleteApplicationsByProjectID(event.AggregateID)
		if err == nil {
			return a.view.ProcessedApplicationSequence(event)
		}
	default:
		return a.view.ProcessedApplicationSequence(event)
	}
	if err != nil {
		return err
	}
	return a.view.PutApplication(app, event)
}

func (a *Application) OnError(event *models.Event, spoolerError error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerError).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerError, a.view.GetLatestApplicationFailedEvent, a.view.ProcessedApplicationFailedEvent, a.view.ProcessedApplicationSequence, a.errorCountUntilSkip)
}

func (a *Application) OnSuccess() error {
	return spooler.HandleSuccess(a.view.UpdateApplicationSpoolerRunTimestamp)
}

func (a *Application) getProjectByID(ctx context.Context, projID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &es_model.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, a.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-DBf32", "Errors.Project.NotFound")
	}

	return es_model.ProjectToModel(esProject), nil
}
