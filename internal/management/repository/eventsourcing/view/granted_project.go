package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	grantedProjectTable = "management.granted_projects"
)

func (v *View) GrantedProjectByIDs(projectID, orgID string) (*model.GrantedProject, error) {
	return view.GrantedProjectByIDs(v.Db, grantedProjectTable, projectID, orgID)
}

func (v *View) SearchGrantedProjects(request *proj_model.GrantedProjectSearchRequest) ([]*model.GrantedProject, int, error) {
	return view.SearchGrantedProjects(v.Db, grantedProjectTable, request)
}

func (v *View) PutGrantedProject(project *model.GrantedProject) error {
	err := view.PutGrantedProject(v.Db, grantedProjectTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedGrantedProjectSequence(project.Sequence)
}

func (v *View) DeleteGrantedProject(projectID, orgID string, eventSequence uint64) error {
	err := view.DeleteGrantedProject(v.Db, grantedProjectTable, projectID, orgID)
	if err != nil {
		return nil
	}
	return v.ProcessedGrantedProjectSequence(eventSequence)
}

func (v *View) GetLatestGrantedProjectSequence() (uint64, error) {
	return v.latestSequence(grantedProjectTable)
}

func (v *View) ProcessedGrantedProjectSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(grantedProjectTable, eventSequence)
}

func (v *View) GetLatestGrantedProjectFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(grantedProjectTable, sequence)
}

func (v *View) ProcessedGrantedProjectFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
