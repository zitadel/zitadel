package query

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

func readModelToIAM(readModel *ReadModel) *model.IAM {
	return &model.IAM{
		ObjectRoot:                      readModelToObjectRoot(readModel.ReadModel),
		GlobalOrgID:                     readModel.GlobalOrgID,
		IAMProjectID:                    readModel.ProjectID,
		SetUpDone:                       readModel.SetUpDone,
		SetUpStarted:                    readModel.SetUpStarted,
		Members:                         readModelToMembers(&readModel.Members),
		DefaultLabelPolicy:              readModelToLabelPolicy(&readModel.DefaultLabelPolicy),
		DefaultLoginPolicy:              readModelToLoginPolicy(&readModel.DefaultLoginPolicy),
		DefaultOrgIAMPolicy:             readModelToOrgIAMPolicy(&readModel.DefaultOrgIAMPolicy),
		DefaultPasswordAgePolicy:        readModelToPasswordAgePolicy(&readModel.DefaultPasswordAgePolicy),
		DefaultPasswordComplexityPolicy: readModelToPasswordComplexityPolicy(&readModel.DefaultPasswordComplexityPolicy),
		DefaultPasswordLockoutPolicy:    readModelToPasswordLockoutPolicy(&readModel.DefaultPasswordLockoutPolicy),
		IDPs:                            readModelToIDPConfigs(&readModel.IDPs),
	}
}

func readModelToIDPConfigView(rm *IAMIDPConfigReadModel) *domain.IDPConfigView {
	converted := &domain.IDPConfigView{
		AggregateID:     rm.AggregateID,
		ChangeDate:      rm.ChangeDate,
		CreationDate:    rm.CreationDate,
		IDPConfigID:     rm.ConfigID,
		IDPProviderType: rm.ProviderType,
		IsOIDC:          rm.OIDCConfig != nil,
		Name:            rm.Name,
		Sequence:        rm.ProcessedSequence,
		State:           rm.State,
		StylingType:     rm.StylingType,
	}
	if rm.OIDCConfig != nil {
		converted.OIDCClientID = rm.OIDCConfig.ClientID
		converted.OIDCClientSecret = rm.OIDCConfig.ClientSecret
		converted.OIDCIDPDisplayNameMapping = rm.OIDCConfig.IDPDisplayNameMapping
		converted.OIDCIssuer = rm.OIDCConfig.Issuer
		converted.OIDCScopes = rm.OIDCConfig.Scopes
		converted.OIDCUsernameMapping = rm.OIDCConfig.UserNameMapping
	}
	return converted
}

func readModelToMember(readModel *MemberReadModel) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: readModelToObjectRoot(readModel.ReadModel),
		Roles:      readModel.Roles,
		UserID:     readModel.UserID,
	}
}

func readModelToMembers(readModel *IAMMembersReadModel) []*model.IAMMember {
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

func readModelToLabelPolicy(readModel *IAMLabelPolicyReadModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.LabelPolicyReadModel.ReadModel),
		PrimaryColor:   readModel.PrimaryColor,
		SecondaryColor: readModel.SecondaryColor,
		Default:        true,
		//TODO: OTPState: int32,
	}
}

func readModelToLoginPolicy(readModel *IAMLoginPolicyReadModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.LoginPolicyReadModel.ReadModel),
		AllowExternalIdp:      readModel.AllowExternalIDP,
		AllowRegister:         readModel.AllowRegister,
		AllowUsernamePassword: readModel.AllowUserNamePassword,
		Default:               true,
		//TODO: IDPProviders: []*model.IDPProvider,
		//TODO: OTPState: int32,
	}
}
func readModelToOrgIAMPolicy(readModel *IAMOrgIAMPolicyReadModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.OrgIAMPolicyReadModel.ReadModel),
		UserLoginMustBeDomain: readModel.UserLoginMustBeDomain,
		Default:               true,
		//TODO: OTPState: int32,
	}
}
func readModelToPasswordAgePolicy(readModel *IAMPasswordAgePolicyReadModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.PasswordAgePolicyReadModel.ReadModel),
		ExpireWarnDays: uint64(readModel.ExpireWarnDays),
		MaxAgeDays:     uint64(readModel.MaxAgeDays),
		//TODO: OTPState: int32,
	}
}
func readModelToPasswordComplexityPolicy(readModel *IAMPasswordComplexityPolicyReadModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   readModelToObjectRoot(readModel.PasswordComplexityPolicyReadModel.ReadModel),
		HasLowercase: readModel.HasLowercase,
		HasNumber:    readModel.HasNumber,
		HasSymbol:    readModel.HasSymbol,
		HasUppercase: readModel.HasUpperCase,
		MinLength:    uint64(readModel.MinLength),
		//TODO: OTPState: int32,
	}
}
func readModelToPasswordLockoutPolicy(readModel *IAMPasswordLockoutPolicyReadModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.PasswordLockoutPolicyReadModel.ReadModel),
		MaxAttempts:         uint64(readModel.MaxAttempts),
		ShowLockOutFailures: readModel.ShowLockOutFailures,
		//TODO: OTPState: int32,
	}
}

func readModelToIDPConfigs(rm *IAMIDPConfigsReadModel) []*model.IDPConfig {
	configs := make([]*model.IDPConfig, len(rm.Configs))
	for i, config := range rm.Configs {
		configs[i] = readModelToIDPConfig(&IAMIDPConfigReadModel{IDPConfigReadModel: *config})
	}
	return configs
}

func readModelToIDPConfig(rm *IAMIDPConfigReadModel) *model.IDPConfig {
	return &model.IDPConfig{
		ObjectRoot:  readModelToObjectRoot(rm.ReadModel),
		OIDCConfig:  readModelToIDPOIDCConfig(rm.OIDCConfig),
		IDPConfigID: rm.ConfigID,
		Name:        rm.Name,
		State:       model.IDPConfigState(rm.State),
		StylingType: model.IDPStylingType(rm.StylingType),
	}
}

func readModelToIDPOIDCConfig(rm *OIDCConfigReadModel) *model.OIDCIDPConfig {
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

func readModelToObjectRoot(readModel eventstore.ReadModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   readModel.AggregateID,
		ChangeDate:    readModel.ChangeDate,
		CreationDate:  readModel.CreationDate,
		ResourceOwner: readModel.ResourceOwner,
		Sequence:      readModel.ProcessedSequence,
	}
}
