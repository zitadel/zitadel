package eventstore

import (
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/project/model"
	proj_view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type ProjectRepo struct {
	View *view.View
}

func (a *ApplicationRepo) ProjectRolesByProjectID(projectID string) ([]*model.ProjectRoleView, error) {
	roles, err := a.View.ProjectRolesByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	return proj_view_model.ProjectRolesToModel(roles), nil
}
