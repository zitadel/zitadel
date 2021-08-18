package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

const (
	applicationTable = "authz.applications"
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

func (k *Application) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (a *Application) ViewModel() string {
	return applicationTable
}

func (p *Application) Subscription() *v1.Subscription {
	return p.subscription
}

func (a *Application) AggregateTypes() []models.AggregateType {
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
	return view.ProjectQuery(sequence.CurrentSequence), nil
}

func (a *Application) Reduce(event *models.Event) (err error) {
	app := new(view_model.ApplicationView)
	switch event.Type {
	case es_model.ApplicationAdded:
		app.AppendEvent(event)
	case es_model.ApplicationChanged,
		es_model.OIDCConfigAdded,
		es_model.OIDCConfigChanged,
		es_model.APIConfigAdded,
		es_model.APIConfigChanged,
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
		return a.view.DeleteApplication(app.ID, event)
	default:
		return a.view.ProcessedApplicationSequence(event)
	}
	if err != nil {
		return err
	}
	return a.view.PutApplication(app, event)
}

func (a *Application) OnError(event *models.Event, spoolerError error) error {
	logging.LogWithFields("SPOOL-sjZw", "id", event.AggregateID).WithError(spoolerError).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerError, a.view.GetLatestApplicationFailedEvent, a.view.ProcessedApplicationFailedEvent, a.view.ProcessedApplicationSequence, a.errorCountUntilSkip)
}

func (a *Application) OnSuccess() error {
	return spooler.HandleSuccess(a.view.UpdateApplicationSpoolerRunTimestamp)
}
