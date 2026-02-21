package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func PasswordExpirySettings() domain.PasswordExpirySettingsRepository {
	return new(passwordExpirySetting)
}

type passwordExpirySetting struct {
	setting[domain.PasswordExpirySetting]
}

func (pes passwordExpirySetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return pes.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.PasswordExpirySettingsRepository].
func (pes passwordExpirySetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordExpirySetting, error) {
	setting, err := pes.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.PasswordExpirySettingsRepository].
func (pes passwordExpirySetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordExpirySetting, error) {
	settings, err := pes.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PasswordExpirySetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	passwordExpirySettingExpireWarnDaysPath = []string{"expireWarnDays"}
	passwordExpirySettingMaxAgeDaysPath     = []string{"maxAgeDays"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.PasswordExpirySettingsRepository].
func (pes passwordExpirySetting) Overwrite(value *domain.PasswordExpirySetting) database.Changes {
	return append(
		pes.setAttributes(&value.Setting),
		pes.SetExpireWarnDays(value.ExpireWarnDays),
		pes.SetMaxAgeDays(value.MaxAgeDays),
	)
}

// SetExpireWarnDays implements [domain.PasswordExpirySettingsRepository].
func (pes passwordExpirySetting) SetExpireWarnDays(value uint64) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordExpirySettingExpireWarnDaysPath, value)
		},
	}
}

// SetMaxAgeDays implements [domain.PasswordExpirySettingsRepository].
func (pes passwordExpirySetting) SetMaxAgeDays(value uint64) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordExpirySettingMaxAgeDaysPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.PasswordExpirySettingsRepository = (*passwordExpirySetting)(nil)
