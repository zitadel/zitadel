package command

import (
	"github.com/zitadel/zitadel/internal/domain"
)

func groupWriteModelToGroup(writeModel *GroupWriteModel) *domain.Group {
	return &domain.Group{
		ObjectRoot:  writeModelToObjectRoot(writeModel.WriteModel),
		Name:        writeModel.Name,
		Description: writeModel.Description,
	}
}

// func roleWriteModelToRole(writeModel *ProjectRoleWriteModel) *domain.ProjectRole {
// 	return &domain.ProjectRole{
// 		ObjectRoot:  writeModelToObjectRoot(writeModel.WriteModel),
// 		Key:         writeModel.Key,
// 		DisplayName: writeModel.DisplayName,
// 		Group:       writeModel.Group,
// 	}
// }

// func memberWriteModelToProjectGrantMember(writeModel *ProjectGrantMemberWriteModel) *domain.ProjectGrantMember {
// 	return &domain.ProjectGrantMember{
// 		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel),
// 		Roles:      writeModel.Roles,
// 		GrantID:    writeModel.GrantID,
// 		UserID:     writeModel.UserID,
// 	}
// }
