package command

import (
	"github.com/caos/zitadel/internal/v2/domain"
)

func projectWriteModelToProject(writeModel *ProjectWriteModel) *domain.Project {
	return &domain.Project{
		ObjectRoot:           writeModelToObjectRoot(writeModel.WriteModel),
		Name:                 writeModel.Name,
		ProjectRoleAssertion: writeModel.ProjectRoleAssertion,
		ProjectRoleCheck:     writeModel.ProjectRoleCheck,
	}
}
