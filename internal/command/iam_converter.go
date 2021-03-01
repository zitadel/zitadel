package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

func writeModelToObjectRoot(writeModel eventstore.WriteModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   writeModel.AggregateID,
		ChangeDate:    writeModel.ChangeDate,
		ResourceOwner: writeModel.ResourceOwner,
		Sequence:      writeModel.ProcessedSequence,
	}
}

func writeModelToIAM(wm *IAMWriteModel) *domain.IAM {
	return &domain.IAM{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		SetUpStarted: wm.SetUpStarted,
		SetUpDone:    wm.SetUpDone,
		GlobalOrgID:  wm.GlobalOrgID,
		IAMProjectID: wm.ProjectID,
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
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel),
		AllowUsernamePassword: wm.AllowUserNamePassword,
		AllowRegister:         wm.AllowRegister,
		AllowExternalIDP:      wm.AllowExternalIDP,
		ForceMFA:              wm.ForceMFA,
		PasswordlessType:      wm.PasswordlessType,
	}
}

func writeModelToLabelPolicy(wm *LabelPolicyWriteModel) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel),
		PrimaryColor:   wm.PrimaryColor,
		SecondaryColor: wm.SecondaryColor,
	}
}

func writeModelToMailTemplate(wm *MailTemplateWriteModel) *domain.MailTemplate {
	return &domain.MailTemplate{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Template:   wm.Template,
	}
}

func writeModelToMailText(wm *MailTextWriteModel) *domain.MailText {
	return &domain.MailText{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		MailTextType: wm.MailTextType,
		Language:     wm.Language,
		Title:        wm.Title,
		PreHeader:    wm.PreHeader,
		Subject:      wm.Subject,
		Greeting:     wm.Greeting,
		Text:         wm.Text,
		ButtonText:   wm.ButtonText,
	}
}

func writeModelToOrgIAMPolicy(wm *IAMOrgIAMPolicyWriteModel) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.PolicyOrgIAMWriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func writeModelToMailTemplatePolicy(wm *MailTemplateWriteModel) *domain.MailTemplate {
	return &domain.MailTemplate{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Template:   wm.Template,
	}
}

func writeModelToMailTextPolicy(wm *MailTextWriteModel) *domain.MailText {
	return &domain.MailText{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		State:        wm.State,
		MailTextType: wm.MailTextType,
		Language:     wm.Language,
		Title:        wm.Title,
		PreHeader:    wm.PreHeader,
		Subject:      wm.Subject,
		Greeting:     wm.Greeting,
		Text:         wm.Text,
		ButtonText:   wm.ButtonText,
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

func writeModelToPasswordLockoutPolicy(wm *PasswordLockoutPolicyWriteModel) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.WriteModel),
		MaxAttempts:         wm.MaxAttempts,
		ShowLockOutFailures: wm.ShowLockOutFailures,
	}
}

func writeModelToIDPConfig(wm *IDPConfigWriteModel) *domain.IDPConfig {
	return &domain.IDPConfig{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		IDPConfigID: wm.ConfigID,
		Name:        wm.Name,
		State:       wm.State,
		StylingType: wm.StylingType,
	}
}

func writeModelToIDPOIDCConfig(wm *OIDCConfigWriteModel) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel),
		ClientID:              wm.ClientID,
		IDPConfigID:           wm.IDPConfigID,
		IDPDisplayNameMapping: wm.IDPDisplayNameMapping,
		Issuer:                wm.Issuer,
		Scopes:                wm.Scopes,
		UsernameMapping:       wm.UserNameMapping,
	}
}

func writeModelToIDPProvider(wm *IdentityProviderWriteModel) *domain.IDPProvider {
	return &domain.IDPProvider{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		IDPConfigID: wm.IDPConfigID,
		Type:        wm.IDPProviderType,
	}
}
