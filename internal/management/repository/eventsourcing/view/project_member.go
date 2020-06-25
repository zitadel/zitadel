package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	projectMemberTable = "management.project_members"
)

func (v *View) ProjectMemberByIDs(projectID, userID string) (*model.ProjectMemberView, error) {
	return view.ProjectMemberByIDs(v.Db, projectMemberTable, projectID, userID)
}

func (v *View) SearchProjectMembers(request *proj_model.ProjectMemberSearchRequest) ([]*model.ProjectMemberView, int, error) {
	return view.SearchProjectMembers(v.Db, projectMemberTable, request)
}

func (v *View) ProjectMembersByUserID(userID string) ([]*model.ProjectMemberView, error) {
	return view.ProjectMembersByUserID(v.Db, projectMemberTable, userID)
}

func (v *View) PutProjectMember(project *model.ProjectMemberView, sequence uint64) error {
	err := view.PutProjectMember(v.Db, projectMemberTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectMemberSequence(sequence)
}

func (v *View) DeleteProjectMember(projectID, userID string, eventSequence uint64) error {
	err := view.DeleteProjectMember(v.Db, projectMemberTable, projectID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectMemberSequence(eventSequence)
}

func (v *View) GetLatestProjectMemberSequence() (uint64, error) {
	return v.latestSequence(projectMemberTable)
}

func (v *View) ProcessedProjectMemberSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(projectMemberTable, eventSequence)
}

func (v *View) GetLatestProjectMemberFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(projectMemberTable, sequence)
}

func (v *View) ProcessedProjectMemberFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
