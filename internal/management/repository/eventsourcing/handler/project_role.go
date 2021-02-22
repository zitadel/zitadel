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
	projectRoleTable = "management.project_roles"
)

type ProjectRole struct {
	handler
	subscription *eventstore.Subscription
}

func newProjectRole(
	handler handler,
) *ProjectRole {
	h := &ProjectRole{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *ProjectRole) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *ProjectRole) ViewModel() string {
	return projectRoleTable
}

func (_ *ProjectRole) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.ProjectAggregate}
}

func (p *ProjectRole) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectRole) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectRoleSequence()
	if err != nil {
		return nil, err
	}
	return proj_view.ProjectQuery(sequence.CurrentSequence), nil
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
		return p.view.DeleteProjectRole(event.AggregateID, event.ResourceOwner, role.Key, event)
	case es_model.ProjectRemoved:
		return p.view.DeleteProjectRolesByProjectID(event.AggregateID)
	default:
		return p.view.ProcessedProjectRoleSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectRole(role, event)
}

func (p *ProjectRole) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-lso9w", "id", event.AggregateID).WithError(err).Warn("something went wrong in project role handler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectRoleFailedEvent, p.view.ProcessedProjectRoleFailedEvent, p.view.ProcessedProjectRoleSequence, p.errorCountUntilSkip)
}

func (p *ProjectRole) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectRoleSpoolerRunTimestamp)
}
