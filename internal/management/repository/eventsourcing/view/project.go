package view

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	projectTable = "management.projects"
)

func (v *View) ProjectByID(projectID string) (*model.ProjectView, error) {
	return view.ProjectByID(v.Db, projectTable, projectID)
}

func (v *View) SearchProjects(request *proj_model.ProjectViewSearchRequest) ([]*model.ProjectView, uint64, error) {
	return view.SearchProjects(v.Db, projectTable, request)
}

func (v *View) PutProject(project *model.ProjectView, event *models.Event) error {
	err := view.PutProject(v.Db, projectTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectSequence(event)
}

func (v *View) DeleteProject(projectID string, event *models.Event) error {
	err := view.DeleteProject(v.Db, projectTable, projectID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectSequence(event)
}

func (v *View) GetLatestProjectSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(projectTable, aggregateType)
}

func (v *View) ProcessedProjectSequence(event *models.Event) error {
	return v.saveCurrentSequence(projectTable, event)
}

func (v *View) UpdateProjectSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(projectTable)
}

func (v *View) GetLatestProjectFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(projectTable, sequence)
}

func (v *View) ProcessedProjectFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
