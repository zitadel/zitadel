package command

import (
	"github.com/caos/zitadel/internal/v2/domain"
)

func orgWriteModelToOrgIAMPolicy(wm *ORGOrgIAMPolicyWriteModel) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.PolicyOrgIAMWriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func orgWriteModelToPasswordComplexityPolicy(wm *OrgPasswordComplexityPolicyWriteModel) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.PasswordComplexityPolicyWriteModel.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUpperCase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}
