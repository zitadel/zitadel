package repository

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func BrandingSettings() domain.BrandingSettingsRepository {
	return new(brandingSetting)
}

type brandingSetting struct {
	setting[domain.BrandingSetting]
}

func (bs brandingSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return bs.ensure(ctx, client, instanceID, organizationID, domain.SettingStatePreview, changes...)
}

// Activate implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return bs.ActivateAt(ctx, client, condition, time.Time{})
}

var (
	brandingSettingActivatePrefix = "INSERT INTO zitadel.settings (instance_id, organization_id, id, type, state, attributes, created_at, updated_at) " +
		"SELECT settings.instance_id, settings.organization_id, settings.id, settings.type, $1::zitadel.settings_state, settings.attributes, settings.created_at, "
	brandingSettingActivateSuffix = " ON CONFLICT (instance_id, organization_id, type, state) DO UPDATE SET (attributes, updated_at) = " +
		"(EXCLUDED.attributes, EXCLUDED.updated_at)"
)

// ActivateAt implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) ActivateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, t time.Time) (int64, error) {
	var updatedAt any = database.NowInstruction
	if !t.IsZero() {
		updatedAt = t
	}
	builder := database.NewStatementBuilder(brandingSettingActivatePrefix, domain.SettingStateActive)
	builder.WriteArg(updatedAt)
	builder.WriteString(" FROM zitadel.settings")
	writeCondition(builder, condition)
	builder.WriteString(brandingSettingActivateSuffix)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// Get implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.BrandingSetting, error) {
	setting, err := bs.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	setting.Attributes.State = setting.State
	return setting.Attributes, nil
}

// List implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.BrandingSetting, error) {
	settings, err := bs.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.BrandingSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
		result[i].State = setting.State
	}
	return result, nil
}

// Set implements [domain.BrandingSettingsRepository].
// The state of the setting will be set to preview, so that it can be activated with [Activate] or [ActivateAt].
func (s brandingSetting) Set(ctx context.Context, client database.QueryExecutor, setting *domain.BrandingSetting) error {
	setting.State = domain.SettingStatePreview
	return s.setting.Set(ctx, client, setting)
}

var (
	brandingSettingPrimaryColorLightPath     = []string{"primaryColorLight"}
	brandingSettingBackgroundColorLightPath  = []string{"backgroundColorLight"}
	brandingSettingWarnColorLightPath        = []string{"warnColorLight"}
	brandingSettingFontColorLightPath        = []string{"fontColorLight"}
	brandingSettingLogoURLLightPath          = []string{"logoUrlLight"}
	brandingSettingIconURLLightPath          = []string{"iconUrlLight"}
	brandingSettingPrimaryColorDarkPath      = []string{"primaryColorDark"}
	brandingSettingBackgroundColorDarkPath   = []string{"backgroundColorDark"}
	brandingSettingWarnColorDarkPath         = []string{"warnColorDark"}
	brandingSettingFontColorDarkPath         = []string{"fontColorDark"}
	brandingSettingLogoURLDarkPath           = []string{"logoUrlDark"}
	brandingSettingIconURLDarkPath           = []string{"iconUrlDark"}
	brandingSettingHideLoginNameSuffixPath   = []string{"hideLoginNameSuffix"}
	brandingSettingErrorMessagePopupPath     = []string{"errorMessagePopup"}
	brandingSettingDisableWatermarkPopupPath = []string{"disableWatermarkPopup"}
	brandingSettingThemeModePath             = []string{"themeMode"}
	brandingSettingFontURLPath               = []string{"fontUrl"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// StateCondition implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) StateCondition(state domain.SettingState) database.Condition {
	return database.NewTextCondition(bs.StateColumn(), database.TextOperationEqual, state.String())
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) Overwrite(value *domain.BrandingSetting) database.Changes {
	return append(bs.setAttributes(&value.Setting),
		bs.setState(value.State),
		bs.SetPrimaryColorLight(value.PrimaryColorLight),
		bs.SetBackgroundColorLight(value.BackgroundColorLight),
		bs.SetWarnColorLight(value.WarnColorLight),
		bs.SetFontColorLight(value.FontColorLight),
		bs.SetLogoURLLight(value.LogoURLLight),
		bs.SetIconURLLight(value.IconURLLight),
		bs.SetPrimaryColorDark(value.PrimaryColorDark),
		bs.SetBackgroundColorDark(value.BackgroundColorDark),
		bs.SetWarnColorDark(value.WarnColorDark),
		bs.SetFontColorDark(value.FontColorDark),
		bs.SetLogoURLDark(value.LogoURLDark),
		bs.SetIconURLDark(value.IconURLDark),
		bs.SetHideLoginNameSuffix(value.HideLoginNameSuffix),
		bs.SetErrorMessagePopup(value.ErrorMessagePopup),
		bs.SetDisableWatermark(value.DisableWatermark),
		bs.SetThemeMode(value.ThemeMode),
		bs.SetFontURL(value.FontURL),
	)
}

// SetBackgroundColorDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetBackgroundColorDark(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingBackgroundColorDarkPath, value)
		},
	}
}

// SetBackgroundColorLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetBackgroundColorLight(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingBackgroundColorLightPath, value)
		},
	}
}

// SetDisableWatermark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetDisableWatermark(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingDisableWatermarkPopupPath, value)
		},
	}
}

// SetErrorMessagePopup implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetErrorMessagePopup(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingErrorMessagePopupPath, value)
		},
	}
}

// SetFontColorDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetFontColorDark(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingFontColorDarkPath, value)
		},
	}
}

// SetFontColorLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetFontColorLight(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingFontColorLightPath, value)
		},
	}
}

// SetFontURL implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetFontURL(value *url.URL) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			if value == nil {
				return database.SetJSONValue[any](column, brandingSettingFontURLPath, nil)
			}
			return database.SetJSONValue(column, brandingSettingFontURLPath, value.String())
		},
	}
}

// SetHideLoginNameSuffix implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetHideLoginNameSuffix(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingHideLoginNameSuffixPath, value)
		},
	}
}

// SetIconURLDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetIconURLDark(value *url.URL) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			if value == nil {
				return database.SetJSONValue[any](column, brandingSettingIconURLDarkPath, nil)
			}
			return database.SetJSONValue(column, brandingSettingIconURLDarkPath, value.String())
		},
	}
}

// SetIconURLLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetIconURLLight(value *url.URL) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			if value == nil {
				return database.SetJSONValue[any](column, brandingSettingIconURLLightPath, nil)
			}
			return database.SetJSONValue(column, brandingSettingIconURLLightPath, value.String())
		},
	}
}

// SetLogoURLDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetLogoURLDark(value *url.URL) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			if value == nil {
				return database.SetJSONValue[any](column, brandingSettingLogoURLDarkPath, nil)
			}
			return database.SetJSONValue(column, brandingSettingLogoURLDarkPath, value.String())
		},
	}
}

// SetLogoURLLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetLogoURLLight(value *url.URL) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			if value == nil {
				return database.SetJSONValue[any](column, brandingSettingLogoURLLightPath, nil)
			}
			return database.SetJSONValue(column, brandingSettingLogoURLLightPath, value.String())
		},
	}
}

// SetPrimaryColorDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetPrimaryColorDark(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingPrimaryColorDarkPath, value)
		},
	}
}

// SetPrimaryColorLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetPrimaryColorLight(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingPrimaryColorLightPath, value)
		},
	}
}

// SetThemeMode implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetThemeMode(value domain.BrandingPolicyThemeMode) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingThemeModePath, value)
		},
	}
}

// SetWarnColorDark implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetWarnColorDark(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingWarnColorDarkPath, value)
		},
	}
}

// SetWarnColorLight implements [domain.BrandingSettingsRepository].
func (bs brandingSetting) SetWarnColorLight(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, brandingSettingWarnColorLightPath, value)
		},
	}
}

var _ domain.BrandingSettingsRepository = (*brandingSetting)(nil)
