package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	projectGrantMemberTable = "management.project_grant_members"
)

func (v *View) ProjectGrantMemberByIDs(projectID, userID string) (*model.ProjectGrantMemberView, error) {
	return view.ProjectGrantMemberByIDs(v.Db, projectGrantMemberTable, projectID, userID)
}

func (v *View) SearchProjectGrantMembers(request *proj_model.ProjectGrantMemberSearchRequest) ([]*model.ProjectGrantMemberView, int, error) {
	return view.SearchProjectGrantMembers(v.Db, projectGrantMemberTable, request)
}

func (v *View) ProjectGrantMembersByUserID(userID string) ([]*model.ProjectGrantMemberView, error) {
	return view.ProjectGrantMembersByUserID(v.Db, projectGrantMemberTable, userID)
}

func (v *View) PutProjectGrantMember(project *model.ProjectGrantMemberView, sequence uint64) error {
	err := view.PutProjectGrantMember(v.Db, projectGrantMemberTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectGrantMemberSequence(sequence)
}

func (v *View) DeleteProjectGrantMember(grantID, userID string, eventSequence uint64) error {
	err := view.DeleteProjectGrantMember(v.Db, projectGrantMemberTable, grantID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectGrantMemberSequence(eventSequence)
}

func (v *View) GetLatestProjectGrantMemberSequence() (uint64, error) {
	return v.latestSequence(projectGrantMemberTable)
}

func (v *View) ProcessedProjectGrantMemberSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(projectGrantMemberTable, eventSequence)
}

func (v *View) GetLatestProjectGrantMemberFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(projectGrantMemberTable, sequence)
}

func (v *View) ProcessedProjectGrantMemberFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
