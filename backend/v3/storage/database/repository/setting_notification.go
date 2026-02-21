package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func NotificationSettings() domain.NotificationSettingsRepository {
	return new(notificationSetting)
}

type notificationSetting struct {
	setting[domain.NotificationSetting]
}

func (ns notificationSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return ns.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.NotificationSettingsRepository].
func (ns notificationSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.NotificationSetting, error) {
	setting, err := ns.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.NotificationSettingsRepository].
func (ns notificationSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.NotificationSetting, error) {
	settings, err := ns.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.NotificationSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	notificationSettingPasswordChangePath = []string{"passwordChange"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.NotificationSettingsRepository].
func (ns notificationSetting) Overwrite(value *domain.NotificationSetting) database.Changes {
	return append(
		ns.setAttributes(&value.Setting),
		ns.SetPasswordChange(value.PasswordChange),
	)
}

// SetPasswordChange implements [domain.NotificationSettingsRepository].
func (ns notificationSetting) SetPasswordChange(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, notificationSettingPasswordChangePath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.NotificationSettingsRepository = (*notificationSetting)(nil)
