package view

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	projectRoleTable = "auth.project_roles"
)

func (v *View) ProjectRoleByIDs(projectID, orgID, key string) (*model.ProjectRoleView, error) {
	return view.ProjectRoleByIDs(v.Db, projectRoleTable, projectID, orgID, key)
}

func (v *View) ProjectRolesByProjectID(projectID string) ([]*model.ProjectRoleView, error) {
	return view.ProjectRolesByProjectID(v.Db, projectRoleTable, projectID)
}

func (v *View) ResourceOwnerProjectRolesByKey(projectID, resourceowner, key string) ([]*model.ProjectRoleView, error) {
	return view.ResourceOwnerProjectRolesByKey(v.Db, projectRoleTable, projectID, resourceowner, key)
}

func (v *View) ResourceOwnerProjectRoles(projectID, resourceowner string) ([]*model.ProjectRoleView, error) {
	return view.ResourceOwnerProjectRoles(v.Db, projectRoleTable, projectID, resourceowner)
}

func (v *View) SearchProjectRoles(request *proj_model.ProjectRoleSearchRequest) ([]*model.ProjectRoleView, uint64, error) {
	return view.SearchProjectRoles(v.Db, projectRoleTable, request)
}

func (v *View) PutProjectRole(role *model.ProjectRoleView, event *models.Event) error {
	err := view.PutProjectRole(v.Db, projectRoleTable, role)
	if err != nil {
		return err
	}
	return v.ProcessedProjectRoleSequence(event)
}

func (v *View) DeleteProjectRole(projectID, orgID, key string, event *models.Event) error {
	err := view.DeleteProjectRole(v.Db, projectRoleTable, projectID, orgID, key)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectRoleSequence(event)
}

func (v *View) DeleteProjectRolesByProjectID(projectID string) error {
	return view.DeleteProjectRolesByProjectID(v.Db, projectRoleTable, projectID)
}

func (v *View) GetLatestProjectRoleSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(projectRoleTable, aggregateType)
}

func (v *View) ProcessedProjectRoleSequence(event *models.Event) error {
	return v.saveCurrentSequence(projectRoleTable, event)
}

func (v *View) UpdateProjectRoleSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(projectRoleTable)
}

func (v *View) GetLatestProjectRoleFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(projectRoleTable, sequence)
}

func (v *View) ProcessedProjectRoleFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
