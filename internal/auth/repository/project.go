package repository

import (
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectRolesByProjectID(projectID string) ([]*model.ProjectRoleView, error)
}
