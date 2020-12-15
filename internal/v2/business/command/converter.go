package command

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/label"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login/idpprovider"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
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

func writeModelToLoginPolicy(wm *login.WriteModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel.WriteModel),
		AllowUsernamePassword: wm.AllowUserNamePassword,
		AllowRegister:         wm.AllowRegister,
		AllowExternalIdp:      wm.AllowExternalIDP,
		ForceMFA:              wm.ForceMFA,
		PasswordlessType:      model.PasswordlessType(wm.PasswordlessType),
	}
}

func writeModelToLabelPolicy(wm *label.WriteModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel.WriteModel),
		PrimaryColor:   wm.PrimaryColor,
		SecondaryColor: wm.SecondaryColor,
	}
}

func writeModelToOrgIAMPolicy(wm *org_iam.WriteModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            writeModelToObjectRoot(wm.WriteModel.WriteModel),
		UserLoginMustBeDomain: wm.UserLoginMustBeDomain,
	}
}

func writeModelToPasswordAgePolicy(wm *password_age.WriteModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel.WriteModel),
		MaxAgeDays:     wm.MaxAgeDays,
		ExpireWarnDays: wm.ExpireWarnDays,
	}
}

func writeModelToPasswordComplexityPolicy(wm *password_complexity.WriteModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   writeModelToObjectRoot(wm.WriteModel.WriteModel),
		MinLength:    wm.MinLength,
		HasLowercase: wm.HasLowercase,
		HasUppercase: wm.HasUpperCase,
		HasNumber:    wm.HasNumber,
		HasSymbol:    wm.HasSymbol,
	}
}

func writeModelToPasswordLockoutPolicy(wm *password_lockout.WriteModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          writeModelToObjectRoot(wm.WriteModel.WriteModel),
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

func writeModelToIDPProvider(wm *idpprovider.WriteModel) *model.IDPProvider {
	return &model.IDPProvider{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel.WriteModel),
		IDPConfigID: wm.IDPConfigID,
		Type:        model.IDPProviderType(wm.IDPProviderType),
	}
}
