package command

import "github.com/caos/zitadel/internal/domain"

func userGrantWriteModelToUserGrant(writeModel *UserGrantWriteModel) *domain.UserGrant {
	return &domain.UserGrant{
		ObjectRoot:     writeModelToObjectRoot(writeModel.WriteModel),
		UserID:         writeModel.UserID,
		ProjectID:      writeModel.ProjectID,
		ProjectGrantID: writeModel.ProjectGrantID,
		RoleKeys:       writeModel.RoleKeys,
		State:          writeModel.State,
	}
}
