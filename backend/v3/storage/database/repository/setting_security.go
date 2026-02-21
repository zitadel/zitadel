package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func SecuritySettings() domain.SecuritySettingsRepository {
	return new(securitySetting)
}

type securitySetting struct {
	setting[domain.SecuritySetting]
}

func (ss securitySetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return ss.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.SecuritySettingsRepository].
func (ss securitySetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecuritySetting, error) {
	setting, err := ss.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.SecuritySettingsRepository].
func (ss securitySetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.SecuritySetting, error) {
	settings, err := ss.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.SecuritySetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	securitySettingEnableIframeEmbeddingPath = []string{"enableIframeEmbedding"}
	securitySettingAllowedOriginsPath        = []string{"allowedOrigins"}
	securitySettingEnableImpersonationPath   = []string{"enableImpersonation"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.SecuritySettingsRepository].
func (ss securitySetting) Overwrite(value *domain.SecuritySetting) database.Changes {
	return append(
		ss.setAttributes(&value.Setting),
		ss.SetEnableIframeEmbedding(value.EnableIframeEmbedding),
		ss.SetAllowedOrigins(value.AllowedOrigins),
		ss.SetEnableImpersonation(value.EnableImpersonation),
	)
}

// SetEnableIframeEmbedding implements [domain.SecuritySettingsRepository].
func (ss securitySetting) SetEnableIframeEmbedding(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, securitySettingEnableIframeEmbeddingPath, value)
		},
	}
}

// SetAllowedOrigins implements [domain.SecuritySettingsRepository].
func (ss securitySetting) SetAllowedOrigins(values []string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, securitySettingAllowedOriginsPath, values)
		},
	}
}

// SetEnableImpersonation implements [domain.SecuritySettingsRepository].
func (ss securitySetting) SetEnableImpersonation(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, securitySettingEnableImpersonationPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.SecuritySettingsRepository = (*securitySetting)(nil)
