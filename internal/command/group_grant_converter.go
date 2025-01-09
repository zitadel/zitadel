package command

import "github.com/zitadel/zitadel/internal/domain"

func groupGrantWriteModelToGroupGrant(writeModel *GroupGrantWriteModel) *domain.GroupGrant {
	return &domain.GroupGrant{
		ObjectRoot:     writeModelToObjectRoot(writeModel.WriteModel),
		GroupID:        writeModel.GroupID,
		ProjectID:      writeModel.ProjectID,
		ProjectGrantID: writeModel.ProjectGrantID,
		RoleKeys:       writeModel.RoleKeys,
		State:          writeModel.State,
	}
}
