package command

import (
	"github.com/caos/zitadel/internal/iam/model"
)

func orgWriteModelToOrgIAMPolicy(wm *ORGOrgIAMPolicyWriteModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.PolicyOrgIAMWriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func orgWriteModelToPasswordComplexityPolicy(wm *OrgPasswordComplexityPolicyWriteModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.PasswordComplexityPolicyWriteModel.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUpperCase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}
