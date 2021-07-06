package query

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
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

func readModelToIDPConfigDomain(rm *IAMIDPConfigReadModel) domain.IDPConfig {
	config := domain.CommonIDPConfig{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  rm.AggregateID,
			ChangeDate:   rm.ChangeDate,
			CreationDate: rm.CreationDate,
			Sequence:     rm.ProcessedSequence,
		},
		IDPConfigID:     rm.ConfigID,
		IDPProviderType: rm.ProviderType,
		//IsOIDC:          rm.OIDCConfig != nil,
		Name:        rm.Name,
		State:       rm.State,
		StylingType: rm.StylingType,
	}
	if rm.OIDCConfig != nil {
		return &domain.OIDCIDPConfig{
			CommonIDPConfig:       		config,
			ClientID:              		rm.OIDCConfig.ClientID,
			ClientSecret:          		rm.OIDCConfig.ClientSecret,
			IDPDisplayNameMapping: 		rm.OIDCConfig.IDPDisplayNameMapping,
			Issuer:                		rm.OIDCConfig.Issuer,
			Scopes:                		rm.OIDCConfig.Scopes,
			UsernameMapping:       		rm.OIDCConfig.UserNameMapping,
			OAuthAuthorizationEndpoint: rm.OIDCConfig.AuthorizationEndpoint,
			OAuthTokenEndpoint: 		rm.OIDCConfig.TokenEndpoint,
		}
	}
	return &domain.AuthConnectorIDPConfig{
		CommonIDPConfig: config,
		BaseURL:         rm.AuthConnectorConfig.BaseURL,
		ProviderID:      rm.AuthConnectorConfig.ProviderID,
		MachineID:       rm.AuthConnectorConfig.MachineID,
	}
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
		ObjectRoot:          readModelToObjectRoot(readModel.LabelPolicyReadModel.ReadModel),
		PrimaryColor:        readModel.PrimaryColor,
		BackgroundColor:     readModel.BackgroundColor,
		WarnColor:           readModel.WarnColor,
		FontColor:           readModel.FontColor,
		PrimaryColorDark:    readModel.PrimaryColorDark,
		BackgroundColorDark: readModel.BackgroundColorDark,
		WarnColorDark:       readModel.WarnColorDark,
		FontColorDark:       readModel.FontColorDark,
		Default:             true,
	}
}

func readModelToLoginPolicy(readModel *IAMLoginPolicyReadModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.LoginPolicyReadModel.ReadModel),
		AllowExternalIdp:      readModel.AllowExternalIDP,
		AllowRegister:         readModel.AllowRegister,
		AllowUsernamePassword: readModel.AllowUserNamePassword,
		Default:               true,
	}
}
func readModelToOrgIAMPolicy(readModel *IAMOrgIAMPolicyReadModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.OrgIAMPolicyReadModel.ReadModel),
		UserLoginMustBeDomain: readModel.UserLoginMustBeDomain,
		Default:               true,
	}
}
func readModelToPasswordAgePolicy(readModel *IAMPasswordAgePolicyReadModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.PasswordAgePolicyReadModel.ReadModel),
		ExpireWarnDays: readModel.ExpireWarnDays,
		MaxAgeDays:     readModel.MaxAgeDays,
	}
}
func readModelToPasswordComplexityPolicy(readModel *IAMPasswordComplexityPolicyReadModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   readModelToObjectRoot(readModel.PasswordComplexityPolicyReadModel.ReadModel),
		HasLowercase: readModel.HasLowercase,
		HasNumber:    readModel.HasNumber,
		HasSymbol:    readModel.HasSymbol,
		HasUppercase: readModel.HasUpperCase,
		MinLength:    readModel.MinLength,
	}
}
func readModelToPasswordLockoutPolicy(readModel *IAMPasswordLockoutPolicyReadModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.PasswordLockoutPolicyReadModel.ReadModel),
		MaxAttempts:         readModel.MaxAttempts,
		ShowLockOutFailures: readModel.ShowLockOutFailures,
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

func readModelToIDPOIDCConfig(rm *IDPOIDCConfigReadModel) *model.OIDCIDPConfig {
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
