package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

const (
	projectTable = "management.projects"
)

type Project struct {
	handler
	subscription *eventstore.Subscription
}

func newProject(handler handler) *Project {
	h := &Project{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *Project) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *Project) ViewModel() string {
	return projectTable
}

func (_ *Project) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.ProjectAggregate}
}

func (p *Project) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *Project) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectSequence()
	if err != nil {
		return nil, err
	}
	return proj_view.ProjectQuery(sequence.CurrentSequence), nil
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
		return p.view.DeleteProject(event.AggregateID, event)
	default:
		return p.view.ProcessedProjectSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProject(project, event)
}

func (p *Project) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-dLsop3", "id", event.AggregateID).WithError(err).Warn("something went wrong in projecthandler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectFailedEvent, p.view.ProcessedProjectFailedEvent, p.view.ProcessedProjectSequence, p.errorCountUntilSkip)
}

func (p *Project) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectSpoolerRunTimestamp)
}
