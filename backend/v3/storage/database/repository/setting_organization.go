package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func OrganizationSettings() domain.OrganizationSettingsRepository {
	return new(organizationSetting)
}

type organizationSetting struct {
	setting[domain.OrganizationSetting]
}

func (os organizationSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return os.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.OrganizationSettingsRepository].
func (os organizationSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationSetting, error) {
	setting, err := os.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.OrganizationSettingsRepository].
func (os organizationSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationSetting, error) {
	settings, err := os.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.OrganizationSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	organizationSettingOrganizationScopedUsernamesPath = []string{"organizationScopedUsernames"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.OrganizationSettingsRepository].
func (os organizationSetting) Overwrite(value *domain.OrganizationSetting) database.Changes {
	return append(
		os.setAttributes(&value.Setting),
		os.SetOrganizationScopedUsernames(value.OrganizationScopedUsernames),
	)
}

// SetOrganizationScopedUsernames implements [domain.OrganizationSettingsRepository].
func (os organizationSetting) SetOrganizationScopedUsernames(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, organizationSettingOrganizationScopedUsernamesPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.OrganizationSettingsRepository = (*organizationSetting)(nil)
