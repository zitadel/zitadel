package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func LoginSettings() domain.LoginSettingsRepository {
	return new(loginSetting)
}

type loginSetting struct {
	setting[domain.LoginSetting]
}

func (ls loginSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return ls.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.LoginSettingsRepository].
func (ls loginSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSetting, error) {
	setting, err := ls.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.LoginSettingsRepository].
func (ls loginSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LoginSetting, error) {
	settings, err := ls.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.LoginSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	loginSettingAllowUsernamePasswordPath       = []string{"allowUsernamePassword"}
	loginSettingAllowRegisterPath               = []string{"allowRegister"}
	loginSettingAllowExternalIdpPath            = []string{"allowExternalIdp"}
	loginSettingForceMultiFactorPath            = []string{"forceMultiFactor"}
	loginSettingForceMultiFactorLocalOnlyPath   = []string{"forceMultiFactorLocalOnly"}
	loginSettingHidePasswordResetPath           = []string{"hidePasswordReset"}
	loginSettingIgnoreUnknownUsernamesPath      = []string{"ignoreUnknownUsernames"}
	loginSettingAllowDomainDiscoveryPath        = []string{"allowDomainDiscovery"}
	loginSettingDisableLoginWithEmailPath       = []string{"disableLoginWithEmail"}
	loginSettingDisableLoginWithPhonePath       = []string{"disableLoginWithPhone"}
	loginSettingPasswordlessTypePath            = []string{"passwordlessType"}
	loginSettingDefaultRedirectUriPath          = []string{"defaultRedirectUri"}
	loginSettingPasswordCheckLifetimePath       = []string{"passwordCheckLifetime"}
	loginSettingExternalLoginCheckLifetimePath  = []string{"externalLoginCheckLifetime"}
	loginSettingMultiFactorInitSkipLifetimePath = []string{"multiFactorInitSkipLifetime"}
	loginSettingSecondFactorCheckLifetimePath   = []string{"secondFactorCheckLifetime"}
	loginSettingMultiFactorCheckLifetimePath    = []string{"multiFactorCheckLifetime"}
	loginSettingMultiFactorTypesPath            = []string{"multiFactorTypes"}
	loginSettingSecondFactorTypesPath           = []string{"secondFactorTypes"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.LoginSettingsRepository].
func (ls loginSetting) Overwrite(value *domain.LoginSetting) database.Changes {
	return append(ls.setAttributes(&value.Setting),
		ls.SetAllowUsernamePassword(value.AllowUsernamePassword),
		ls.SetAllowRegister(value.AllowRegister),
		ls.SetAllowExternalIDP(value.AllowExternalIDP),
		ls.SetForceMultiFactor(value.ForceMultiFactor),
		ls.SetForceMultiFactorLocalOnly(value.ForceMultiFactorLocalOnly),
		ls.SetHidePasswordReset(value.HidePasswordReset),
		ls.SetIgnoreUnknownUsernames(value.IgnoreUnknownUsernames),
		ls.SetAllowDomainDiscovery(value.AllowDomainDiscovery),
		ls.SetDisableLoginWithEmail(value.DisableLoginWithEmail),
		ls.SetDisableLoginWithPhone(value.DisableLoginWithPhone),
		ls.SetPasswordlessType(value.PasswordlessType),
		ls.SetDefaultRedirectURI(value.DefaultRedirectURI),
		ls.SetPasswordCheckLifetime(value.PasswordCheckLifetime),
		ls.SetExternalLoginCheckLifetime(value.ExternalLoginCheckLifetime),
		ls.SetMultiFactorInitSkipLifetime(value.MultiFactorInitSkipLifetime),
		ls.SetSecondFactorCheckLifetime(value.SecondFactorCheckLifetime),
		ls.SetMultiFactorCheckLifetime(value.MultiFactorCheckLifetime),
		ls.setMultiFactorTypes(value.MultiFactorTypes),
		ls.setSecondFactorTypes(value.SecondFactorTypes),
	)
}

// AddMultiFactorType implements [domain.LoginSettingsRepository].
func (ls loginSetting) AddMultiFactorType(value domain.MultiFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingMultiFactorTypesPath, value)
		},
	}
}

// RemoveMultiFactorType implements [domain.LoginSettingsRepository].
func (ls loginSetting) RemoveMultiFactorType(value domain.MultiFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.RemoveFromJSONArray(column, loginSettingMultiFactorTypesPath, value)
		},
	}
}

func (ls loginSetting) setMultiFactorTypes(value []domain.MultiFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingMultiFactorTypesPath, value)
		},
	}
}

// AddSecondFactorType implements [domain.LoginSettingsRepository].
func (ls loginSetting) AddSecondFactorType(value domain.SecondFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.AppendToJSONArray(column, loginSettingSecondFactorTypesPath, value)
		},
	}
}

// RemoveSecondFactorType implements [domain.LoginSettingsRepository].
func (ls loginSetting) RemoveSecondFactorType(value domain.SecondFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.RemoveFromJSONArray(column, loginSettingSecondFactorTypesPath, value)
		},
	}
}

func (ls loginSetting) setSecondFactorTypes(value []domain.SecondFactorType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingSecondFactorTypesPath, value)
		},
	}
}

// SetAllowDomainDiscovery implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetAllowDomainDiscovery(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingAllowDomainDiscoveryPath, value)
		},
	}
}

// SetAllowExternalIDP implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetAllowExternalIDP(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingAllowExternalIdpPath, value)
		},
	}
}

// SetAllowRegister implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetAllowRegister(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingAllowRegisterPath, value)
		},
	}
}

// SetAllowUsernamePassword implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetAllowUsernamePassword(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingAllowUsernamePasswordPath, value)
		},
	}
}

// SetDefaultRedirectURI implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetDefaultRedirectURI(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingDefaultRedirectUriPath, value)
		},
	}
}

// SetDisableLoginWithEmail implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetDisableLoginWithEmail(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingDisableLoginWithEmailPath, value)
		},
	}
}

// SetDisableLoginWithPhone implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetDisableLoginWithPhone(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingDisableLoginWithPhonePath, value)
		},
	}
}

// SetExternalLoginCheckLifetime implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetExternalLoginCheckLifetime(value time.Duration) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingExternalLoginCheckLifetimePath, value)
		},
	}
}

// SetForceMultiFactor implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetForceMultiFactor(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingForceMultiFactorPath, value)
		},
	}
}

// SetForceMultiFactorLocalOnly implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetForceMultiFactorLocalOnly(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingForceMultiFactorLocalOnlyPath, value)
		},
	}
}

// SetHidePasswordReset implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetHidePasswordReset(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingHidePasswordResetPath, value)
		},
	}
}

// SetIgnoreUnknownUsernames implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetIgnoreUnknownUsernames(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingIgnoreUnknownUsernamesPath, value)
		},
	}
}

// SetMultiFactorInitSkipLifetime implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetMultiFactorInitSkipLifetime(value time.Duration) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingMultiFactorInitSkipLifetimePath, value)
		},
	}
}

// SetMultiFactorCheckLifetime implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetMultiFactorCheckLifetime(value time.Duration) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingMultiFactorCheckLifetimePath, value)
		},
	}
}

// SetPasswordCheckLifetime implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetPasswordCheckLifetime(value time.Duration) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingPasswordCheckLifetimePath, value)
		},
	}
}

// SetPasswordlessType implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetPasswordlessType(value domain.PasswordlessType) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingPasswordlessTypePath, value)
		},
	}
}

// SetSecondFactorCheckLifetime implements [domain.LoginSettingsRepository].
func (ls loginSetting) SetSecondFactorCheckLifetime(value time.Duration) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, loginSettingSecondFactorCheckLifetimePath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.LoginSettingsRepository = (*loginSetting)(nil)
