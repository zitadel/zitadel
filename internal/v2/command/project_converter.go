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

func projectGrantWriteModelToProjectGrant(writeModel *ProjectGrantWriteModel) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot:   writeModelToObjectRoot(writeModel.WriteModel),
		GrantID:      writeModel.GrantID,
		GrantedOrgID: writeModel.GrantedOrgID,
		RoleKeys:     writeModel.RoleKeys,
		State:        writeModel.State,
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

func roleWriteModelToRole(writeModel *ProjectRoleWriteModel) *domain.ProjectRole {
	return &domain.ProjectRole{
		ObjectRoot:  writeModelToObjectRoot(writeModel.WriteModel),
		Key:         writeModel.Key,
		DisplayName: writeModel.DisplayName,
		Group:       writeModel.Group,
	}
}

func memberWriteModelToProjectGrantMember(writeModel *ProjectGrantMemberWriteModel) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel),
		Roles:      writeModel.Roles,
		GrantID:    writeModel.GrantID,
		UserID:     writeModel.UserID,
	}
}
