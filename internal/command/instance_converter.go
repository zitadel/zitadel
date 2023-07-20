package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

func writeModelToObjectRoot(writeModel eventstore.WriteModel) models.ObjectRoot {
	return models.ObjectRoot{
		InstanceID:    writeModel.InstanceID,
		AggregateID:   writeModel.AggregateID,
		ChangeDate:    writeModel.ChangeDate,
		ResourceOwner: writeModel.ResourceOwner,
		Sequence:      writeModel.ProcessedSequence,
	}
}

func memberWriteModelToMember(writeModel *MemberWriteModel) *domain.Member {
	return &domain.Member{
		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel),
		Roles:      writeModel.Roles,
		UserID:     writeModel.UserID,
	}
}

func writeModelToLoginPolicy(wm *LoginPolicyWriteModel) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot:                 writeModelToObjectRoot(wm.WriteModel),
		AllowUsernamePassword:      wm.AllowUserNamePassword,
		AllowRegister:              wm.AllowRegister,
		AllowExternalIDP:           wm.AllowExternalIDP,
		HidePasswordReset:          wm.HidePasswordReset,
		IgnoreUnknownUsernames:     wm.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       wm.AllowDomainDiscovery,
		ForceMFA:                   wm.ForceMFA,
		ForceMFALocalOnly:          wm.ForceMFALocalOnly,
		PasswordlessType:           wm.PasswordlessType,
		DefaultRedirectURI:         wm.DefaultRedirectURI,
		PasswordCheckLifetime:      wm.PasswordCheckLifetime,
		ExternalLoginCheckLifetime: wm.ExternalLoginCheckLifetime,
		MFAInitSkipLifetime:        wm.MFAInitSkipLifetime,
		SecondFactorCheckLifetime:  wm.SecondFactorCheckLifetime,
		MultiFactorCheckLifetime:   wm.MultiFactorCheckLifetime,
	}
}

func writeModelToLabelPolicy(wm *LabelPolicyWriteModel) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.WriteModel),
		PrimaryColor:        wm.PrimaryColor,
		BackgroundColor:     wm.BackgroundColor,
		WarnColor:           wm.WarnColor,
		FontColor:           wm.FontColor,
		PrimaryColorDark:    wm.PrimaryColorDark,
		BackgroundColorDark: wm.BackgroundColorDark,
		WarnColorDark:       wm.WarnColorDark,
		FontColorDark:       wm.FontColorDark,
		HideLoginNameSuffix: wm.HideLoginNameSuffix,
		ErrorMsgPopup:       wm.ErrorMsgPopup,
		DisableWatermark:    wm.DisableWatermark,
	}
}

func writeModelToMailTemplate(wm *MailTemplateWriteModel) *domain.MailTemplate {
	return &domain.MailTemplate{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Template:   wm.Template,
	}
}

func writeModelToDomainPolicy(wm *InstanceDomainPolicyWriteModel) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		ObjectRoot:                             writeModelToObjectRoot(wm.PolicyDomainWriteModel.WriteModel),
		UserLoginMustBeDomain:                  wm.UserLoginMustBeDomain,
		ValidateOrgDomains:                     wm.ValidateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: wm.SMTPSenderAddressMatchesInstanceDomain,
	}
}

func writeModelToMailTemplatePolicy(wm *MailTemplateWriteModel) *domain.MailTemplate {
	return &domain.MailTemplate{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Template:   wm.Template,
	}
}

func writeModelToCustomText(wm *CustomTextWriteModel) *domain.CustomText {
	return &domain.CustomText{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		State:      wm.State,
		Key:        wm.Key,
		Language:   wm.Language,
		Text:       wm.Text,
	}
}

func writeModelToPasswordAgePolicy(wm *PasswordAgePolicyWriteModel) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel),
		MaxAgeDays:     wm.MaxAgeDays,
		ExpireWarnDays: wm.ExpireWarnDays,
	}
}

func writeModelToPasswordComplexityPolicy(wm *PasswordComplexityPolicyWriteModel) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUppercase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}

func writeModelToLockoutPolicy(wm *LockoutPolicyWriteModel) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.WriteModel),
		MaxPasswordAttempts: wm.MaxPasswordAttempts,
		ShowLockOutFailures: wm.ShowLockOutFailures,
	}
}

func writeModelToPrivacyPolicy(wm *PrivacyPolicyWriteModel) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		TOSLink:      wm.TOSLink,
		PrivacyLink:  wm.PrivacyLink,
		HelpLink:     wm.HelpLink,
		SupportEmail: wm.SupportEmail,
	}
}

func writeModelToIDPConfig(wm *IDPConfigWriteModel) *domain.IDPConfig {
	return &domain.IDPConfig{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		IDPConfigID:  wm.ConfigID,
		Name:         wm.Name,
		State:        wm.State,
		StylingType:  wm.StylingType,
		AutoRegister: wm.AutoRegister,
	}
}

func writeModelToIDPOIDCConfig(wm *OIDCConfigWriteModel) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel),
		ClientID:              wm.ClientID,
		IDPConfigID:           wm.IDPConfigID,
		IDPDisplayNameMapping: wm.IDPDisplayNameMapping,
		Issuer:                wm.Issuer,
		AuthorizationEndpoint: wm.AuthorizationEndpoint,
		TokenEndpoint:         wm.TokenEndpoint,
		Scopes:                wm.Scopes,
		UsernameMapping:       wm.UserNameMapping,
	}
}

func writeModelToIDPJWTConfig(wm *JWTConfigWriteModel) *domain.JWTIDPConfig {
	return &domain.JWTIDPConfig{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		IDPConfigID:  wm.IDPConfigID,
		JWTEndpoint:  wm.JWTEndpoint,
		Issuer:       wm.Issuer,
		KeysEndpoint: wm.KeysEndpoint,
		HeaderName:   wm.HeaderName,
	}
}

func writeModelToIDPProvider(wm *IdentityProviderWriteModel) *domain.IDPProvider {
	return &domain.IDPProvider{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		IDPConfigID: wm.IDPConfigID,
		Type:        wm.IDPProviderType,
	}
}
