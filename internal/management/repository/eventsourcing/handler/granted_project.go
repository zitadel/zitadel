package handler

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type GrantedProject struct {
	handler
	eventstore    eventstore.Eventstore
	projectEvents *proj_event.ProjectEventstore
}

const (
	grantedProjectTable = "management.granted_projects"
)

func (p *GrantedProject) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *GrantedProject) ViewModel() string {
	return grantedProjectTable
}

func (p *GrantedProject) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestGrantedProjectSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.ProjectQuery(sequence), nil
}

func (p *GrantedProject) Process(event *models.Event) (err error) {
	grantedProject := new(view_model.GrantedProjectView)
	switch event.Type {
	case es_model.ProjectAdded:
		grantedProject.AppendEvent(event)
	case es_model.ProjectChanged:
		grantedProject, err = p.view.GrantedProjectByIDs(event.AggregateID, event.ResourceOwner)
		if err != nil {
			return err
		}
		err = grantedProject.AppendEvent(event)
		if err != nil {
			return err
		}
		p.updateExistingProjects(grantedProject)
	case es_model.ProjectDeactivated, es_model.ProjectReactivated:
		grantedProject, err = p.view.GrantedProjectByIDs(event.AggregateID, event.ResourceOwner)
		if err != nil {
			return err
		}
		err = grantedProject.AppendEvent(event)
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
		//TODO: read org
	case es_model.ProjectGrantChanged:
		grant := new(view_model.ProjectGrant)
		err := grant.SetData(event)
		if err != nil {
			return err
		}
		grantedProject, err = p.view.GrantedProjectByIDs(event.AggregateID, grant.GrantedOrgID)
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
		return p.view.DeleteGrantedProject(event.AggregateID, grant.GrantedOrgID, event.Sequence)
	default:
		return p.view.ProcessedGrantedProjectSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutGrantedProject(grantedProject)
}

func (p *GrantedProject) getOrg(orgID string) {
	//TODO: Get Org
}

func (p *GrantedProject) getProject(projectID string) (*model.Project, error) {
	return p.projectEvents.ProjectByID(context.Background(), projectID)
}

func (p *GrantedProject) updateExistingProjects(project *view_model.GrantedProjectView) {
	projects, err := p.view.GrantedProjectsByID(project.ProjectID)
	if err != nil {
		logging.LogWithFields("SPOOL-los03", "id", project.ProjectID).WithError(err).Warn("could not update existing projects")
	}
	for _, existing := range projects {
		existing.Name = project.Name
		err := p.view.PutGrantedProject(existing)
		logging.LogWithFields("SPOOL-sjwi3", "id", existing.ProjectID).OnError(err).Warn("could not update existing project")
	}
}

func (p *GrantedProject) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-is8wa", "id", event.AggregateID).WithError(err).Warn("something went wrong in granted projecthandler")
	return spooler.HandleError(event, err, p.view.GetLatestGrantedProjectFailedEvent, p.view.ProcessedGrantedProjectFailedEvent, p.view.ProcessedGrantedProjectSequence, p.errorCountUntilSkip)
}
