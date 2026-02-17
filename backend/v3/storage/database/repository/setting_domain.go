package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func DomainSettings() domain.DomainSettingsRepository {
	return new(domainSetting)
}

type domainSetting struct {
	setting[domain.DomainSetting]
}

func (ds domainSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return ds.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.DomainSettingsRepository].
func (ds domainSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.DomainSetting, error) {
	setting, err := ds.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.DomainSettingsRepository].
func (ds domainSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.DomainSetting, error) {
	settings, err := ds.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.DomainSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	domainSettingLoginNameIncludesDomainPath                = []string{"loginNameIncludesDomain"}
	domainSettingRequireOrgDomainVerificationPath           = []string{"requireOrgDomainVerification"}
	domainSettingSMTPSenderAddressMatchesInstanceDomainPath = []string{"smtpSenderAddressMatchesInstanceDomain"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.DomainSettingsRepository].
func (ds domainSetting) Overwrite(value *domain.DomainSetting) database.Changes {
	return append(
		ds.setAttributes(&value.Setting),
		ds.SetLoginNameIncludesDomain(value.LoginNameIncludesDomain),
		ds.SetRequireOrgDomainVerification(value.RequireOrgDomainVerification),
		ds.SetSMTPSenderAddressMatchesInstanceDomain(value.SMTPSenderAddressMatchesInstanceDomain),
	)
}

// SetLoginNameIncludesDomain implements [domain.DomainSettingsRepository].
func (ds domainSetting) SetLoginNameIncludesDomain(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, domainSettingLoginNameIncludesDomainPath, value)
		},
	}
}

// SetRequireOrgDomainVerification implements [domain.DomainSettingsRepository].
func (ds domainSetting) SetRequireOrgDomainVerification(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, domainSettingRequireOrgDomainVerificationPath, value)
		},
	}
}

// SetSMTPSenderAddressMatchesInstanceDomain implements [domain.DomainSettingsRepository].
func (ds domainSetting) SetSMTPSenderAddressMatchesInstanceDomain(value bool) database.Change {
	return &attributeChange{
		change: func(column database.Column) database.Change {
			return database.SetJSONValue(column, domainSettingSMTPSenderAddressMatchesInstanceDomainPath, value)
		},
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.DomainSettingsRepository = (*domainSetting)(nil)
