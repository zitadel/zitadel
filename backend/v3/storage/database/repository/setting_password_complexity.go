package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func PasswordComplexitySettings() domain.PasswordComplexitySettingsRepository {
	return new(passwordComplexitySetting)
}

type passwordComplexitySetting struct {
	setting[domain.PasswordComplexitySetting]
}

func (pcs passwordComplexitySetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return pcs.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordComplexitySetting, error) {
	setting, err := pcs.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordComplexitySetting, error) {
	settings, err := pcs.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PasswordComplexitySetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	passwordComplexitySettingMinLengthPath    = []string{"minLength"}
	passwordComplexitySettingHasLowercasePath = []string{"hasLowercase"}
	passwordComplexitySettingHasUppercasePath = []string{"hasUppercase"}
	passwordComplexitySettingHasNumberPath    = []string{"hasNumber"}
	passwordComplexitySettingHasSymbolPath    = []string{"hasSymbol"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) Overwrite(value *domain.PasswordComplexitySetting) database.Changes {
	return append(
		pcs.setAttributes(&value.Setting),
		pcs.SetHasLowercase(value.HasLowercase),
		pcs.SetHasNumber(value.HasNumber),
		pcs.SetHasSymbol(value.HasSymbol),
		pcs.SetHasUppercase(value.HasUppercase),
		pcs.SetMinLength(value.MinLength),
	)
}

// SetHasLowercase implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) SetHasLowercase(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordComplexitySettingHasLowercasePath, value)
		},
	}
}

// SetHasNumber implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) SetHasNumber(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordComplexitySettingHasNumberPath, value)
		},
	}
}

// SetHasSymbol implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) SetHasSymbol(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordComplexitySettingHasSymbolPath, value)
		},
	}
}

// SetHasUppercase implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) SetHasUppercase(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordComplexitySettingHasUppercasePath, value)
		},
	}
}

// SetMinLength implements [domain.PasswordComplexitySettingsRepository].
func (pcs passwordComplexitySetting) SetMinLength(value uint64) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, passwordComplexitySettingMinLengthPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.PasswordComplexitySettingsRepository = (*passwordComplexitySetting)(nil)
