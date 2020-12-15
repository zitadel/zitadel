package query

import (
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/label"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func readModelToIAM(readModel *ReadModel) *model.IAM {
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
		IDPs:                            readModelToIDPConfigs(&readModel.IDPs),
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

func readModelToLabelPolicy(readModel *label.ReadModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.ReadModel.ReadModel),
		PrimaryColor:   readModel.PrimaryColor,
		SecondaryColor: readModel.SecondaryColor,
		Default:        true,
		//TODO: State: int32,
	}
}

func readModelToLoginPolicy(readModel *login.ReadModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.ReadModel.ReadModel),
		AllowExternalIdp:      readModel.AllowExternalIDP,
		AllowRegister:         readModel.AllowRegister,
		AllowUsernamePassword: readModel.AllowUserNamePassword,
		Default:               true,
		//TODO: IDPProviders: []*model.IDPProvider,
		//TODO: State: int32,
	}
}
func readModelToOrgIAMPolicy(readModel *org_iam.ReadModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.ReadModel.ReadModel),
		UserLoginMustBeDomain: readModel.UserLoginMustBeDomain,
		Default:               true,
		//TODO: State: int32,
	}
}
func readModelToPasswordAgePolicy(readModel *password_age.ReadModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.ReadModel.ReadModel),
		ExpireWarnDays: uint64(readModel.ExpireWarnDays),
		MaxAgeDays:     uint64(readModel.MaxAgeDays),
		//TODO: State: int32,
	}
}
func readModelToPasswordComplexityPolicy(readModel *password_complexity.ReadModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   readModelToObjectRoot(readModel.ReadModel.ReadModel),
		HasLowercase: readModel.HasLowercase,
		HasNumber:    readModel.HasNumber,
		HasSymbol:    readModel.HasSymbol,
		HasUppercase: readModel.HasUpperCase,
		MinLength:    uint64(readModel.MinLength),
		//TODO: State: int32,
	}
}
func readModelToPasswordLockoutPolicy(readModel *password_lockout.ReadModel) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.ReadModel.ReadModel),
		MaxAttempts:         uint64(readModel.MaxAttempts),
		ShowLockOutFailures: readModel.ShowLockOutFailures,
		//TODO: State: int32,
	}
}

func readModelToIDPConfigs(rm *iam.IDPConfigsReadModel) []*model.IDPConfig {
	configs := make([]*model.IDPConfig, len(rm.Configs))
	for i, config := range rm.Configs {
		configs[i] = readModelToIDPConfig(&iam.IDPConfigReadModel{ConfigReadModel: *config})
	}
	return configs
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

func readModelToObjectRoot(readModel eventstore.ReadModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   readModel.AggregateID,
		ChangeDate:    readModel.ChangeDate,
		CreationDate:  readModel.CreationDate,
		ResourceOwner: readModel.ResourceOwner,
		Sequence:      readModel.ProcessedSequence,
	}
}
