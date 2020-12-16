package command

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func writeModelToObjectRoot(writeModel eventstore.WriteModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   writeModel.AggregateID,
		ChangeDate:    writeModel.ChangeDate,
		ResourceOwner: writeModel.ResourceOwner,
		Sequence:      writeModel.ProcessedSequence,
	}
}

func writeModelToMember(writeModel *IAMMemberWriteModel) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: writeModelToObjectRoot(writeModel.MemberWriteModel.WriteModel),
		Roles:      writeModel.Roles,
		UserID:     writeModel.UserID,
	}
}

func writeModelToLoginPolicy(wm *IAMLoginPolicyWriteModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.LoginPolicyWriteModel.WriteModel),
		AllowUsernamePassword: wm.AllowUserNamePassword,
		AllowRegister:         wm.AllowRegister,
		AllowExternalIdp:      wm.AllowExternalIDP,
		ForceMFA:              wm.ForceMFA,
		PasswordlessType:      model.PasswordlessType(wm.PasswordlessType),
	}
}

func writeModelToLabelPolicy(wm *IAMLabelPolicyWriteModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.LabelPolicyWriteModel.WriteModel),
		PrimaryColor:   wm.PrimaryColor,
		SecondaryColor: wm.SecondaryColor,
	}
}

func writeModelToOrgIAMPolicy(wm *IAMOrgIAMPolicyWriteModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.PolicyOrgIAMWriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func writeModelToPasswordAgePolicy(wm *IAMPasswordAgePolicyWriteModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.PasswordAgePolicyWriteModel.WriteModel),
		MaxAgeDays:     wm.MaxAgeDays,
		ExpireWarnDays: wm.ExpireWarnDays,
	}
}

func writeModelToPasswordComplexityPolicy(wm *IAMPasswordComplexityPolicyWriteModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.PasswordComplexityPolicyWriteModel.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUpperCase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}

func writeModelToPasswordLockoutPolicy(wm *IAMPasswordLockoutPolicyWriteModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.PasswordLockoutPolicyWriteModel.WriteModel),
		MaxAttempts:         wm.MaxAttempts,
		ShowLockOutFailures: wm.ShowLockOutFailures,
	}
}

func writeModelToIDPConfig(wm *iam.IDPConfigWriteModel) *model.IDPConfig {
	return &model.IDPConfig{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		OIDCConfig:  writeModelToIDPOIDCConfig(wm.OIDCConfig),
		IDPConfigID: wm.ConfigID,
		Name:        wm.Name,
		State:       model.IDPConfigState(wm.State),
		StylingType: model.IDPStylingType(wm.StylingType),
	}
}

func writeModelToIDPOIDCConfig(wm *oidc.ConfigWriteModel) *model.OIDCIDPConfig {
	return &model.OIDCIDPConfig{
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel),
		ClientID:              wm.ClientID,
		IDPConfigID:           wm.IDPConfigID,
		IDPDisplayNameMapping: model.OIDCMappingField(wm.IDPDisplayNameMapping),
		Issuer:                wm.Issuer,
		Scopes:                wm.Scopes,
		UsernameMapping:       model.OIDCMappingField(wm.UserNameMapping),
	}
}

func writeModelToIDPProvider(wm *IAMIdentityProviderWriteModel) *model.IDPProvider {
	return &model.IDPProvider{
		ObjectRoot:  writeModelToObjectRoot(wm.IdentityProviderWriteModel.WriteModel),
		IDPConfigID: wm.IDPConfigID,
		Type:        model.IDPProviderType(wm.IDPProviderType),
	}
}
