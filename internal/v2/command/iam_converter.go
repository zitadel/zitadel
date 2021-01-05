package command

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

func writeModelToObjectRoot(writeModel eventstore.WriteModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   writeModel.AggregateID,
		ChangeDate:    writeModel.ChangeDate,
		ResourceOwner: writeModel.ResourceOwner,
		Sequence:      writeModel.ProcessedSequence,
	}
}

func writeModelToIAM(wm *IAMWriteModel) *model.IAM {
	return &model.IAM{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel),
		SetUpStarted: wm.SetUpStarted,
		SetUpDone:    wm.SetUpDone,
		GlobalOrgID:  wm.GlobalOrgID,
		IAMProjectID: wm.ProjectID,
	}
}

func writeModelToMember(writeModel *IAMMemberWriteModel) *domain.IAMMember {
	return &domain.IAMMember{
		ObjectRoot: writeModelToObjectRoot(writeModel.MemberWriteModel.WriteModel),
		Roles:      writeModel.Roles,
		UserID:     writeModel.UserID,
	}
}

func writeModelToLoginPolicy(wm *IAMLoginPolicyWriteModel) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.LoginPolicyWriteModel.WriteModel),
		AllowUsernamePassword: wm.AllowUserNamePassword,
		AllowRegister:         wm.AllowRegister,
		AllowExternalIdp:      wm.AllowExternalIDP,
		ForceMFA:              wm.ForceMFA,
		PasswordlessType:      wm.PasswordlessType,
	}
}

func writeModelToLabelPolicy(wm *IAMLabelPolicyWriteModel) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.LabelPolicyWriteModel.WriteModel),
		PrimaryColor:   wm.PrimaryColor,
		SecondaryColor: wm.SecondaryColor,
	}
}

func writeModelToOrgIAMPolicy(wm *IAMOrgIAMPolicyWriteModel) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.PolicyOrgIAMWriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func writeModelToPasswordAgePolicy(wm *IAMPasswordAgePolicyWriteModel) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.PasswordAgePolicyWriteModel.WriteModel),
		MaxAgeDays:     wm.MaxAgeDays,
		ExpireWarnDays: wm.ExpireWarnDays,
	}
}

func writeModelToPasswordComplexityPolicy(wm *IAMPasswordComplexityPolicyWriteModel) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.PasswordComplexityPolicyWriteModel.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUpperCase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}

func writeModelToPasswordLockoutPolicy(wm *IAMPasswordLockoutPolicyWriteModel) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.PasswordLockoutPolicyWriteModel.WriteModel),
		MaxAttempts:         wm.MaxAttempts,
		ShowLockOutFailures: wm.ShowLockOutFailures,
	}
}

func writeModelToIDPConfig(wm *IAMIDPConfigWriteModel) *domain.IDPConfig {
	return &domain.IDPConfig{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		OIDCConfig:  writeModelToIDPOIDCConfig(wm.OIDCConfig),
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

func writeModelToIDPProvider(wm *IAMIdentityProviderWriteModel) *domain.IDPProvider {
	return &domain.IDPProvider{
		ObjectRoot:  writeModelToObjectRoot(wm.IdentityProviderWriteModel.WriteModel),
		IDPConfigID: wm.IDPConfigID,
		Type:        wm.IDPProviderType,
	}
}
