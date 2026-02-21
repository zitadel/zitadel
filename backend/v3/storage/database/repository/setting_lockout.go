package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func LockoutSettings() domain.LockoutSettingsRepository {
	return new(lockoutSetting)
}

type lockoutSetting struct {
	setting[domain.LockoutSetting]
}

func (ls lockoutSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return ls.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LockoutSetting, error) {
	setting, err := ls.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LockoutSetting, error) {
	settings, err := ls.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.LockoutSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	lockoutSettingMaxPasswordAttemptsPath = []string{"maxPasswordAttempts"}
	lockoutSettingMaxOTPAttemptsPath      = []string{"maxOtpAttempts"}
	lockoutSettingShowLockOutFailuresPath = []string{"showLockOutFailures"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) Overwrite(value *domain.LockoutSetting) database.Changes {
	return append(
		ls.setAttributes(&value.Setting),
		ls.SetMaxPasswordAttempts(value.MaxPasswordAttempts),
		ls.SetMaxOTPAttempts(value.MaxOTPAttempts),
		ls.SetShowLockOutFailures(value.ShowLockOutFailures),
	)
}

// SetMaxPasswordAttempts implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) SetMaxPasswordAttempts(value uint64) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, lockoutSettingMaxPasswordAttemptsPath, value)
		},
	}
}

// SetMaxOTPAttempts implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) SetMaxOTPAttempts(value uint64) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, lockoutSettingMaxOTPAttemptsPath, value)
		},
	}
}

// SetShowLockOutFailures implements [domain.LockoutSettingsRepository].
func (ls lockoutSetting) SetShowLockOutFailures(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, lockoutSettingShowLockOutFailuresPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.LockoutSettingsRepository = (*lockoutSetting)(nil)
