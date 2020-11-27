package iam

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

func readModelToIAM(readModel *iam_repo.ReadModel) *model.IAM {
	return &model.IAM{
		ObjectRoot:                      readModelToObjectRoot(readModel.ReadModel),
		GlobalOrgID:                     readModel.GlobalOrgID,
		IAMProjectID:                    readModel.ProjectID,
		SetUpDone:                       model.Step(readModel.SetUpDone),
		SetUpStarted:                    model.Step(readModel.SetUpStarted),
		Members:                         readModelToMembers(&readModel.Members),
		DefaultLabelPolicy:              readModelToLabelPolicy(&readModel.DefaultLabelPolicy),
		DefaultLoginPolicy:              readModelToLoginPolicy(&readModel.DefaultLoginPolicy),
		DefaultOrgIAMPolicy:             readModelToOrgIAMPolicy(&readModel.DefaultOrgIAMPolicy),
		DefaultPasswordAgePolicy:        readModelToPasswordAgePolicy(&readModel.DefaultPasswordAgePolicy),
		DefaultPasswordComplexityPolicy: readModelToPasswordComplexityPolicy(&readModel.DefaultPasswordComplexityPolicy),
		DefaultPasswordLockoutPolicy:    readModelToPasswordLockoutPolicy(&readModel.DefaultPasswordLockoutPolicy),
		// TODO: IDPs: []*model.IDPConfig,
	}
}

func readModelToMembers(readModel *iam_repo.MembersReadModel) []*model.IAMMember {
	members := make([]*model.IAMMember, len(readModel.Members))

	for i, member := range readModel.Members {
		members[i] = &model.IAMMember{
			ObjectRoot: readModelToObjectRoot(member.ReadModel),
			Roles:      member.Roles,
			UserID:     member.UserID,
		}
	}

	return members
}

func readModelToLabelPolicy(readModel *iam.LabelPolicyReadModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.ReadModel),
		PrimaryColor:   readModel.PrimaryColor,
		SecondaryColor: readModel.SecondaryColor,
		Default:        true,
		//TODO: State: int32,
	}
}

func readModelToLoginPolicy(readModel *iam.LoginPolicyReadModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.ReadModel),
		AllowExternalIdp:      readModel.AllowExternalIDP,
		AllowRegister:         readModel.AllowRegister,
		AllowUsernamePassword: readModel.AllowUserNamePassword,
		Default:               true,
		//TODO: IDPProviders: []*model.IDPProvider,
		//TODO: State: int32,
	}
}
func readModelToOrgIAMPolicy(readModel *iam.OrgIAMPolicyReadModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.ReadModel),
		UserLoginMustBeDomain: readModel.UserLoginMustBeDomain,
		Default:               true,
		//TODO: State: int32,
	}
}
func readModelToPasswordAgePolicy(readModel *iam.PasswordAgePolicyReadModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.ReadModel),
		ExpireWarnDays: uint64(readModel.ExpireWarnDays),
		MaxAgeDays:     uint64(readModel.MaxAgeDays),
		//TODO: State: int32,
	}
}
func readModelToPasswordComplexityPolicy(readModel *iam.PasswordComplexityPolicyReadModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   readModelToObjectRoot(readModel.ReadModel),
		HasLowercase: readModel.HasLowercase,
		HasNumber:    readModel.HasNumber,
		HasSymbol:    readModel.HasSymbol,
		HasUppercase: readModel.HasUpperCase,
		MinLength:    uint64(readModel.MinLength),
		//TODO: State: int32,
	}
}
func readModelToPasswordLockoutPolicy(readModel *iam.PasswordLockoutPolicyReadModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.ReadModel),
		MaxAttempts:         uint64(readModel.MaxAttempts),
		ShowLockOutFailures: readModel.ShowLockOutFailures,
		//TODO: State: int32,
	}
}

func readModelToObjectRoot(readModel eventstore.ReadModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   readModel.AggregateID,
		ChangeDate:    readModel.ChangeDate,
		CreationDate:  readModel.CreationDate,
		ResourceOwner: readModel.ResourceOwner,
		Sequence:      readModel.ProcessedSequence,
	}
}

func writeModelToObjectRoot(writeModel eventstore.WriteModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   writeModel.AggregateID,
		ChangeDate:    writeModel.ChangeDate,
		ResourceOwner: writeModel.ResourceOwner,
		Sequence:      writeModel.ProcessedSequence,
	}
}

func readModelToMember(readModel *member.ReadModel) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: readModelToObjectRoot(readModel.ReadModel),
		Roles:      readModel.Roles,
		UserID:     readModel.UserID,
	}
}

func writeModelToMember(writeModel *iam.MemberWriteModel) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel.WriteModel),
		Roles:      writeModel.Roles,
		UserID:     writeModel.UserID,
	}
}

func readModelToIDPConfigView(rm *iam.IDPConfigReadModel) *model.IDPConfigView {
	return &model.IDPConfigView{
		AggregateID:               rm.AggregateID,
		ChangeDate:                rm.ChangeDate,
		CreationDate:              rm.CreationDate,
		IDPConfigID:               rm.ConfigID,
		IDPProviderType:           model.IDPProviderType(rm.ProviderType),
		IsOIDC:                    rm.OIDCConfig != nil,
		Name:                      rm.Name,
		OIDCClientID:              rm.OIDCConfig.ClientID,
		OIDCClientSecret:          rm.OIDCConfig.ClientSecret,
		OIDCIDPDisplayNameMapping: model.OIDCMappingField(rm.OIDCConfig.IDPDisplayNameMapping),
		OIDCIssuer:                rm.OIDCConfig.Issuer,
		OIDCScopes:                rm.OIDCConfig.Scopes,
		OIDCUsernameMapping:       model.OIDCMappingField(rm.OIDCConfig.UserNameMapping),
		Sequence:                  rm.ProcessedSequence,
		State:                     model.IDPConfigState(rm.State),
		StylingType:               model.IDPStylingType(rm.StylingType),
	}
}

func readModelToIDPConfig(rm *iam.IDPConfigReadModel) *model.IDPConfig {
	return &model.IDPConfig{
		ObjectRoot:  readModelToObjectRoot(rm.ReadModel),
		OIDCConfig:  readModelToIDPOIDCConfig(rm.OIDCConfig),
		IDPConfigID: rm.ConfigID,
		Name:        rm.Name,
		State:       model.IDPConfigState(rm.State),
		StylingType: model.IDPStylingType(rm.StylingType),
	}
}

func readModelToIDPOIDCConfig(rm *oidc.ConfigReadModel) *model.OIDCIDPConfig {
	return &model.OIDCIDPConfig{
		ObjectRoot:            readModelToObjectRoot(rm.ReadModel),
		ClientID:              rm.ClientID,
		ClientSecret:          rm.ClientSecret,
		ClientSecretString:    string(rm.ClientSecret.Crypted),
		IDPConfigID:           rm.IDPConfigID,
		IDPDisplayNameMapping: model.OIDCMappingField(rm.IDPDisplayNameMapping),
		Issuer:                rm.Issuer,
		Scopes:                rm.Scopes,
		UsernameMapping:       model.OIDCMappingField(rm.UserNameMapping),
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
