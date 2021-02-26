package command

import (
	"github.com/caos/zitadel/internal/domain"
)

func orgWriteModelToOrg(wm *OrgWriteModel) *domain.Org {
	return &domain.Org{
		ObjectRoot:    writeModelToObjectRoot(wm.WriteModel),
		Name:          wm.Name,
		State:         wm.State,
		PrimaryDomain: wm.PrimaryDomain,
	}
}

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
		HasUppercase: wm.HasUppercase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}

func orgDomainWriteModelToOrgDomain(wm *OrgDomainWriteModel) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel),
		Domain:         wm.Domain,
		Primary:        wm.Primary,
		Verified:       wm.Verified,
		ValidationType: wm.ValidationType,
		ValidationCode: wm.ValidationCode,
	}
}
