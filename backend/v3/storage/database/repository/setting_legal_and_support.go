package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func LegalAndSupportSettings() domain.LegalAndSupportSettingsRepository {
	return new(legalAndSupportSetting)
}

type legalAndSupportSetting struct {
	setting[domain.LegalAndSupportSetting]
}

func (lss legalAndSupportSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return lss.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LegalAndSupportSetting, error) {
	setting, err := lss.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LegalAndSupportSetting, error) {
	settings, err := lss.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.LegalAndSupportSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	legalAndSupportSettingTOSLinkPath           = []string{"tosLink"}
	legalAndSupportSettingPrivacyPolicyLinkPath = []string{"privacyPolicyLink"}
	legalAndSupportSettingHelpLinkPath          = []string{"helpLink"}
	legalAndSupportSettingSupportEmailPath      = []string{"supportEmail"}
	legalAndSupportSettingDocsLinkPath          = []string{"docsLink"}
	legalAndSupportSettingCustomLinkPath        = []string{"customLink"}
	legalAndSupportSettingCustomLinkTextPath    = []string{"customLinkText"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) Overwrite(value *domain.LegalAndSupportSetting) database.Changes {
	return append(
		lss.setAttributes(&value.Setting),
		lss.SetTOSLink(value.TOSLink),
		lss.SetPrivacyPolicyLink(value.PrivacyPolicyLink),
		lss.SetHelpLink(value.HelpLink),
		lss.SetSupportEmail(value.SupportEmail),
		lss.SetDocsLink(value.DocsLink),
		lss.SetCustomLink(value.CustomLink),
		lss.SetCustomLinkText(value.CustomLinkText),
	)
}

// SetTOSLink implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetTOSLink(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingTOSLinkPath, value)
		},
	}
}

// SetPrivacyPolicyLink implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetPrivacyPolicyLink(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingPrivacyPolicyLinkPath, value)
		},
	}
}

// SetHelpLink implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetHelpLink(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingHelpLinkPath, value)
		},
	}
}

// SetSupportEmail implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetSupportEmail(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingSupportEmailPath, value)
		},
	}
}

// SetDocsLink implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetDocsLink(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingDocsLinkPath, value)
		},
	}
}

// SetCustomLink implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetCustomLink(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingCustomLinkPath, value)
		},
	}
}

// SetCustomLinkText implements [domain.LegalAndSupportSettingsRepository].
func (lss legalAndSupportSetting) SetCustomLinkText(value string) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, legalAndSupportSettingCustomLinkTextPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.LegalAndSupportSettingsRepository = (*legalAndSupportSetting)(nil)
