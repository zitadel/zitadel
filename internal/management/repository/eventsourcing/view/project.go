package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	projectTable = "management.projects"
)

func (v *View) ProjectByID(projectID string) (*model.ProjectView, error) {
	return view.ProjectByID(v.Db, projectTable, projectID)
}

func (v *View) SearchProjects(request *proj_model.ProjectViewSearchRequest) ([]*model.ProjectView, int, error) {
	return view.SearchProjects(v.Db, projectTable, request)
}

func (v *View) PutProject(project *model.ProjectView) error {
	err := view.PutProject(v.Db, projectTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectSequence(project.Sequence)
}

func (v *View) DeleteProject(projectID string, eventSequence uint64) error {
	err := view.DeleteProject(v.Db, projectTable, projectID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectSequence(eventSequence)
}

func (v *View) GetLatestProjectSequence() (uint64, error) {
	return v.latestSequence(projectTable)
}

func (v *View) ProcessedProjectSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(projectTable, eventSequence)
}

func (v *View) GetLatestProjectFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(projectTable, sequence)
}

func (v *View) ProcessedProjectFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
