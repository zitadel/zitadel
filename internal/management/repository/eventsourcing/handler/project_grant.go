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
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

const (
	grantedProjectTable = "management.project_grants"
)

type ProjectGrant struct {
	handler
	subscription *v1.Subscription
}

func newProjectGrant(
	handler handler,
) *ProjectGrant {
	h := &ProjectGrant{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *ProjectGrant) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *ProjectGrant) ViewModel() string {
	return grantedProjectTable
}

func (p *ProjectGrant) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *ProjectGrant) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.ProjectAggregate}
}

func (p *ProjectGrant) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestProjectGrantSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *ProjectGrant) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectGrantSequence()
	if err != nil {
		return nil, err
	}
	return proj_view.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *ProjectGrant) Reduce(event *models.Event) (err error) {
	grantedProject := new(view_model.ProjectGrantView)
	switch event.Type {
	case es_model.ProjectChanged:
		project, err := p.view.ProjectByID(event.AggregateID)
		if err != nil {
			return err
		}
		return p.updateExistingProjects(project, event)
	case es_model.ProjectGrantAdded:
		err = grantedProject.AppendEvent(event)
		if err != nil {
			return err
		}
		project, err := p.getProject(grantedProject.ProjectID)
		if err != nil {
			return err
		}
		grantedProject.Name = project.Name

		org, err := p.getOrgByID(context.TODO(), grantedProject.OrgID)
		if err != nil {
			return err
		}
		resourceOwner, err := p.getOrgByID(context.TODO(), grantedProject.ResourceOwner)
		if err != nil {
			return err
		}
		p.fillOrgData(grantedProject, org, resourceOwner)
	case es_model.ProjectGrantChanged, es_model.ProjectGrantCascadeChanged:
		grant := new(view_model.ProjectGrant)
		err = grant.SetData(event)
		if err != nil {
			return err
		}
		grantedProject, err = p.view.ProjectGrantByID(grant.GrantID)
		if err != nil {
			return err
		}
		err = grantedProject.AppendEvent(event)
	case es_model.ProjectGrantRemoved:
		grant := new(view_model.ProjectGrant)
		err := grant.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectGrant(grant.GrantID, event)
	case es_model.ProjectRemoved:
		err = p.view.DeleteProjectGrantsByProjectID(event.AggregateID)
		if err != nil {
			return err
		}
		return p.view.ProcessedProjectGrantSequence(event)
	default:
		return p.view.ProcessedProjectGrantSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrant(grantedProject, event)
}

func (p *ProjectGrant) fillOrgData(grantedProject *view_model.ProjectGrantView, org, resourceOwner *org_model.Org) {
	grantedProject.OrgName = org.Name
	grantedProject.ResourceOwnerName = resourceOwner.Name
}

func (p *ProjectGrant) getProject(projectID string) (*proj_model.Project, error) {
	return p.getProjectByID(context.Background(), projectID)
}

func (p *ProjectGrant) updateExistingProjects(project *view_model.ProjectView, event *models.Event) error {
	projectGrants, err := p.view.ProjectGrantsByProjectID(project.ProjectID)
	if err != nil {
		logging.LogWithFields("SPOOL-los03", "id", project.ProjectID).WithError(err).Warn("could not update existing projects")
	}
	for _, existingGrant := range projectGrants {
		existingGrant.Name = project.Name
	}
	return p.view.PutProjectGrants(projectGrants, event)
}

func (p *ProjectGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-sQqOg", "id", event.AggregateID).WithError(err).Warn("something went wrong in granted projecthandler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectGrantFailedEvent, p.view.ProcessedProjectGrantFailedEvent, p.view.ProcessedProjectGrantSequence, p.errorCountUntilSkip)
}

func (p *ProjectGrant) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProjectGrantSpoolerRunTimestamp)
}

func (u *ProjectGrant) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *ProjectGrant) getProjectByID(ctx context.Context, projID string) (*proj_model.Project, error) {
	query, err := proj_view.ProjectByIDQuery(projID, 0)
	if err != nil {
		return nil, err
	}
	esProject := &es_model.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esProject.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esProject.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-NBrw2", "Errors.Project.NotFound")
	}

	return es_model.ProjectToModel(esProject), nil
}
