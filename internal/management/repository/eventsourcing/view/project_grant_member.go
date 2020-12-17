package view

import (
	"github.com/caos/zitadel/internal/eventstore/models"
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

func (v *View) ProjectGrantMembersByProjectID(projectID string) ([]*model.ProjectGrantMemberView, error) {
	return view.ProjectGrantMembersByProjectID(v.Db, projectGrantMemberTable, projectID)
}

func (v *View) SearchProjectGrantMembers(request *proj_model.ProjectGrantMemberSearchRequest) ([]*model.ProjectGrantMemberView, uint64, error) {
	return view.SearchProjectGrantMembers(v.Db, projectGrantMemberTable, request)
}

func (v *View) ProjectGrantMembersByUserID(userID string) ([]*model.ProjectGrantMemberView, error) {
	return view.ProjectGrantMembersByUserID(v.Db, projectGrantMemberTable, userID)
}

func (v *View) PutProjectGrantMember(member *model.ProjectGrantMemberView, event *models.Event) error {
	err := view.PutProjectGrantMember(v.Db, projectGrantMemberTable, member)
	if err != nil {
		return err
	}
	return v.ProcessedProjectGrantMemberSequence(event)
}

func (v *View) PutProjectGrantMembers(members []*model.ProjectGrantMemberView, event *models.Event) error {
	err := view.PutProjectGrantMembers(v.Db, projectGrantMemberTable, members...)
	if err != nil {
		return err
	}
	return v.ProcessedProjectGrantMemberSequence(event)
}

func (v *View) DeleteProjectGrantMember(grantID, userID string, event *models.Event) error {
	err := view.DeleteProjectGrantMember(v.Db, projectGrantMemberTable, grantID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectGrantMemberSequence(event)
}

func (v *View) DeleteProjectGrantMembersByProjectID(projectID string) error {
	return view.DeleteProjectGrantMembersByProjectID(v.Db, projectGrantMemberTable, projectID)
}

func (v *View) GetLatestProjectGrantMemberSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(projectGrantMemberTable, aggregateType)
}

func (v *View) ProcessedProjectGrantMemberSequence(event *models.Event) error {
	return v.saveCurrentSequence(projectGrantMemberTable, event)
}

func (v *View) UpdateProjectGrantMemberSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(projectGrantMemberTable)
}

func (v *View) GetLatestProjectGrantMemberFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(projectGrantMemberTable, sequence)
}

func (v *View) ProcessedProjectGrantMemberFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
