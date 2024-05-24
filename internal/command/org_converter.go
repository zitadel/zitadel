package command

import (
	"github.com/zitadel/zitadel/internal/domain"
)

func orgWriteModelToOrg(wm *OrgWriteModel) *domain.Org {
	return &domain.Org{
		ObjectRoot:    writeModelToObjectRoot(wm.WriteModel),
		Name:          wm.Name,
		State:         wm.State,
		PrimaryDomain: wm.PrimaryDomain,
	}
}

func orgWriteModelToDomainPolicy(wm *OrgDomainPolicyWriteModel) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		ObjectRoot:                             writeModelToObjectRoot(wm.PolicyDomainWriteModel.WriteModel),
		UserLoginMustBeDomain:                  wm.UserLoginMustBeDomain,
		ValidateOrgDomains:                     wm.ValidateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: wm.SMTPSenderAddressMatchesInstanceDomain,
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

func orgWriteModelToPrivacyPolicy(wm *OrgPrivacyPolicyWriteModel) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.PrivacyPolicyWriteModel.WriteModel),
		TOSLink:        wm.TOSLink,
		PrivacyLink:    wm.PrivacyLink,
		HelpLink:       wm.HelpLink,
		SupportEmail:   wm.SupportEmail,
		DocsLink:       wm.DocsLink,
		CustomLink:     wm.CustomLink,
		CustomLinkText: wm.CustomLinkText,
	}
}
