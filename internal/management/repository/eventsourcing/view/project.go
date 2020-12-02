package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
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

func (v *View) PutProject(project *model.ProjectView, eventTimestamp time.Time) error {
	err := view.PutProject(v.Db, projectTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectSequence(project.Sequence, eventTimestamp)
}

func (v *View) DeleteProject(projectID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteProject(v.Db, projectTable, projectID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestProjectSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(projectTable)
}

func (v *View) ProcessedProjectSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(projectTable, eventSequence, eventTimestamp)
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
