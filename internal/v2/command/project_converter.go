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

func applicationWriteModelToApplication(writeModel *ApplicationWriteModel) *domain.Application {
	return &domain.Application{
		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel),
		AppID:      writeModel.AggregateID,
		State:      writeModel.State,
		Name:       writeModel.Name,
		Type:       writeModel.Type,
		//TODO: OIDC Config
	}
}
