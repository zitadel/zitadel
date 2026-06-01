package repository_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestGetLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LoginSettingsRepository()

	settings := []*domain.LoginSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LoginSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LoginSettingsRepository()

	settings := []*domain.LoginSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LoginSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.LoginSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LoginSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LoginSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LoginSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettingsRepository()

	existingSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri2"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LoginSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettingsRepository()

	existingInstanceSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(instanceID, gu.Ptr("foo"), domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettingsRepository()

	settings := []*domain.BrandingSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.Activate(t.Context(), tx, repo.UniqueCondition(settings[0].InstanceID, settings[0].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview))
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.Activate(t.Context(), tx, repo.UniqueCondition(settings[3].InstanceID, settings[3].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview))
	require.NoError(t, err)
	activeOrganizationSetting := *settings[3]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.BrandingSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name: "ok, instance",
			condition: database.And(
				repo.PrimaryKeyCondition(settings[0].InstanceID, settings[0].ID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[0],
		},
		{
			name: "ok, instance, active",
			condition: database.And(
				repo.PrimaryKeyCondition(activeInstanceSetting.InstanceID, activeInstanceSetting.ID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: &activeInstanceSetting,
		},
		{
			name: "ok, organization",
			condition: database.And(
				repo.PrimaryKeyCondition(settings[3].InstanceID, settings[3].ID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[3],
		},
		{
			name: "ok, organization",
			condition: database.And(
				repo.PrimaryKeyCondition(activeOrganizationSetting.InstanceID, activeOrganizationSetting.ID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: &activeOrganizationSetting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettingsRepository()

	settings := []*domain.BrandingSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.Activate(t.Context(), tx, repo.UniqueCondition(settings[0].InstanceID, settings[0].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview))
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.Activate(t.Context(), tx, repo.UniqueCondition(settings[2].InstanceID, settings[2].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview))
	require.NoError(t, err)
	activeOrganizationSetting := *settings[2]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.BrandingSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.BrandingSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.BrandingSettings{&activeOrganizationSetting, settings[2], &activeInstanceSetting, settings[0]},
		},
		{
			name: "all from instance, preview",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: []*domain.BrandingSettings{settings[2], settings[0]},
		},
		{
			name: "all from instance, active",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: []*domain.BrandingSettings{&activeOrganizationSetting, &activeInstanceSetting},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.BrandingSettings{&activeInstanceSetting, settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.BrandingSettings{&activeOrganizationSetting, settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn(), repo.StateColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Len(t, got, len(tt.want)) {
				return
			}
			for i := range got {
				assert.EqualExportedValues(t, tt.want[i], got[i])
			}
		})
	}
}

func TestSetBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.BrandingSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestActivateBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingInstanceSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		wantRowsAffected int64
		condition        database.Condition
		wantErr          error
	}{
		{
			name: "activate instance",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 1,
		},
		{
			name: "activate organization",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				existingOrganizationSettings.OrganizationID,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 1,
		},
		{
			name: "non-existing instance",
			condition: repo.UniqueCondition(
				"foo",
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 0,
		},

		{
			name: "non-existing organization",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				gu.Ptr("foo"),
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := repo.Activate(t.Context(), savepoint, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestActivateAtBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingInstanceSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		wantRowsAffected int64
		condition        database.Condition
		updatedAt        time.Time
		wantErr          error
	}{
		{
			name: "activate instance",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			updatedAt:        time.Now(),
			wantRowsAffected: 1,
		},
		{
			name: "activate organization",
			condition: repo.UniqueCondition(
				existingOrganizationSettings.InstanceID,
				existingOrganizationSettings.OrganizationID,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			updatedAt:        time.Now(),
			wantRowsAffected: 1,
		},
		{
			name: "non-existing instance",
			condition: repo.UniqueCondition(
				"foo",
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			updatedAt:        time.Now(),
			wantRowsAffected: 0,
		},

		{
			name: "non-existing organization",
			condition: repo.UniqueCondition(
				existingOrganizationSettings.InstanceID,
				gu.Ptr("foo"),
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			updatedAt:        time.Now(),
			wantRowsAffected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := repo.ActivateAt(t.Context(), savepoint, tt.condition, tt.updatedAt)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestDeleteBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingInstanceSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx,
		repo.UniqueCondition(
			existingInstanceSettings.InstanceID,
			existingInstanceSettings.OrganizationID,
			domain.SettingTypeBranding,
			domain.SettingStatePreview,
		),
	)
	require.NoError(t, err)
	activeInstanceSetting := *existingInstanceSettings
	activeInstanceSetting.State = domain.SettingStateActive

	existingOrganizationSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx,
		repo.UniqueCondition(
			existingOrganizationSettings.InstanceID,
			existingOrganizationSettings.OrganizationID,
			domain.SettingTypeBranding,
			domain.SettingStatePreview,
		),
	)
	require.NoError(t, err)
	activeOrganizationSetting := *existingOrganizationSettings
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance, active",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization, active",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice, active",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	settings := []*domain.PasswordComplexitySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(8)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(12)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(14)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordComplexitySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	settings := []*domain.PasswordComplexitySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(8)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(12)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordComplexitySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.PasswordComplexitySettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordComplexitySettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordComplexitySettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordComplexitySettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	existingSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.PasswordComplexitySettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(10)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeletePasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	existingInstanceSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypePasswordComplexity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettingsRepository()

	settings := []*domain.PasswordExpirySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(10)),
				MaxAgeDays:     gu.Ptr(uint64(30)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(11)),
				MaxAgeDays:     gu.Ptr(uint64(31)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(12)),
				MaxAgeDays:     gu.Ptr(uint64(32)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(13)),
				MaxAgeDays:     gu.Ptr(uint64(33)),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordExpirySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettingsRepository()

	settings := []*domain.PasswordExpirySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(10)),
				MaxAgeDays:     gu.Ptr(uint64(30)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(11)),
				MaxAgeDays:     gu.Ptr(uint64(31)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(12)),
				MaxAgeDays:     gu.Ptr(uint64(32)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(13)),
				MaxAgeDays:     gu.Ptr(uint64(33)),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordExpirySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.PasswordExpirySettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordExpirySettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordExpirySettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordExpirySettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordExpirySettingsRepository()

	existingSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.PasswordExpirySettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeletePasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordExpirySettingsRepository()

	existingInstanceSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypePasswordExpiry, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettingsRepository()

	settings := []*domain.LockoutSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(3)),
				MaxOTPAttempts:      gu.Ptr(uint64(3)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(4)),
				MaxOTPAttempts:      gu.Ptr(uint64(4)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(5)),
				MaxOTPAttempts:      gu.Ptr(uint64(5)),
				ShowLockOutFailures: gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(6)),
				MaxOTPAttempts:      gu.Ptr(uint64(6)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LockoutSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettingsRepository()

	settings := []*domain.LockoutSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(3)),
				MaxOTPAttempts:      gu.Ptr(uint64(3)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(4)),
				MaxOTPAttempts:      gu.Ptr(uint64(4)),
				ShowLockOutFailures: gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(5)),
				MaxOTPAttempts:      gu.Ptr(uint64(5)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(6)),
				MaxOTPAttempts:      gu.Ptr(uint64(6)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LockoutSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.LockoutSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LockoutSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LockoutSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LockoutSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LockoutSettingsRepository()

	existingSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(3)),
			MaxOTPAttempts:      gu.Ptr(uint64(3)),
			ShowLockOutFailures: gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LockoutSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(4)),
					MaxOTPAttempts:      gu.Ptr(uint64(4)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(5)),
					MaxOTPAttempts:      gu.Ptr(uint64(5)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(3)),
					MaxOTPAttempts:      gu.Ptr(uint64(3)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(3)),
					MaxOTPAttempts:      gu.Ptr(uint64(3)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LockoutSettingsRepository()

	existingInstanceSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(3)),
			MaxOTPAttempts:      gu.Ptr(uint64(3)),
			ShowLockOutFailures: gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(4)),
			MaxOTPAttempts:      gu.Ptr(uint64(4)),
			ShowLockOutFailures: gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeLockout, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLockout, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLockout, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLockout, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLockout, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecuritySettingsRepository()

	settings := []*domain.SecuritySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin1", "origin2"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin3", "origin4"},
				EnableImpersonation:   gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin5", "origin6"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(false),
				AllowedOrigins:        []string{"origin7", "origin8"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.SecuritySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecuritySettingsRepository()

	settings := []*domain.SecuritySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin1", "origin2"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin3", "origin4"},
				EnableImpersonation:   gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(true),
				AllowedOrigins:        []string{"origin5", "origin6"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
				EnableIframeEmbedding: gu.Ptr(false),
				AllowedOrigins:        []string{"origin7", "origin8"},
				EnableImpersonation:   gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.SecuritySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.SecuritySettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.SecuritySettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.SecuritySettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.SecuritySettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecuritySettingsRepository()

	existingSettings := &domain.SecuritySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
			EnableIframeEmbedding: gu.Ptr(true),
			AllowedOrigins:        []string{"origin1", "origin2"},
			EnableImpersonation:   gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.SecuritySettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.SecuritySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
					EnableIframeEmbedding: gu.Ptr(true),
					AllowedOrigins:        []string{"origin1", "origin2"},
					EnableImpersonation:   gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.SecuritySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
					EnableIframeEmbedding: gu.Ptr(true),
					AllowedOrigins:        []string{"origin1", "origin2"},
					EnableImpersonation:   gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.SecuritySettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
					EnableIframeEmbedding: gu.Ptr(true),
					AllowedOrigins:        []string{"origin1", "origin2"},
					EnableImpersonation:   gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.SecuritySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
					EnableIframeEmbedding: gu.Ptr(true),
					AllowedOrigins:        []string{"origin1", "origin2"},
					EnableImpersonation:   gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecuritySettingsRepository()

	existingInstanceSettings := &domain.SecuritySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
			EnableIframeEmbedding: gu.Ptr(true),
			AllowedOrigins:        []string{"origin1", "origin2"},
			EnableImpersonation:   gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.SecuritySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		SecuritySettingsAttributes: domain.SecuritySettingsAttributes{
			EnableIframeEmbedding: gu.Ptr(true),
			AllowedOrigins:        []string{"origin3", "origin4"},
			EnableImpersonation:   gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeSecurity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecurity, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecurity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecurity, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecurity, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.DomainSettingsRepository()

	settings := []*domain.DomainSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(false),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(false),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(false),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.DomainSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.DomainSettingsRepository()

	settings := []*domain.DomainSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(false),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(false),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			DomainSettingsAttributes: domain.DomainSettingsAttributes{
				LoginNameIncludesDomain:                gu.Ptr(true),
				RequireOrgDomainVerification:           gu.Ptr(true),
				SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(false),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.DomainSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.DomainSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.DomainSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.DomainSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.DomainSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.DomainSettingsRepository()

	existingSettings := &domain.DomainSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		DomainSettingsAttributes: domain.DomainSettingsAttributes{
			LoginNameIncludesDomain:                gu.Ptr(true),
			RequireOrgDomainVerification:           gu.Ptr(true),
			SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.DomainSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.DomainSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				DomainSettingsAttributes: domain.DomainSettingsAttributes{
					LoginNameIncludesDomain:                gu.Ptr(true),
					RequireOrgDomainVerification:           gu.Ptr(true),
					SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.DomainSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				DomainSettingsAttributes: domain.DomainSettingsAttributes{
					LoginNameIncludesDomain:                gu.Ptr(true),
					RequireOrgDomainVerification:           gu.Ptr(true),
					SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.DomainSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				DomainSettingsAttributes: domain.DomainSettingsAttributes{
					LoginNameIncludesDomain:                gu.Ptr(true),
					RequireOrgDomainVerification:           gu.Ptr(true),
					SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.DomainSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				DomainSettingsAttributes: domain.DomainSettingsAttributes{
					LoginNameIncludesDomain:                gu.Ptr(true),
					RequireOrgDomainVerification:           gu.Ptr(true),
					SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.DomainSettingsRepository()

	existingInstanceSettings := &domain.DomainSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		DomainSettingsAttributes: domain.DomainSettingsAttributes{
			LoginNameIncludesDomain:                gu.Ptr(true),
			RequireOrgDomainVerification:           gu.Ptr(true),
			SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.DomainSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		DomainSettingsAttributes: domain.DomainSettingsAttributes{
			LoginNameIncludesDomain:                gu.Ptr(true),
			RequireOrgDomainVerification:           gu.Ptr(false),
			SMTPSenderAddressMatchesInstanceDomain: gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeDomain, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeDomain, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeDomain, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeDomain, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeDomain, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.OrganizationSettingsRepository()

	settings := []*domain.OrganizationSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.OrganizationSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.OrganizationSettingsRepository()

	settings := []*domain.OrganizationSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
				OrganizationScopedUsernames: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.OrganizationSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.OrganizationSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.OrganizationSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.OrganizationSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.OrganizationSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.OrganizationSettingsRepository()

	existingSettings := &domain.OrganizationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
			OrganizationScopedUsernames: gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.OrganizationSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.OrganizationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.OrganizationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.OrganizationSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.OrganizationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
					OrganizationScopedUsernames: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.OrganizationSettingsRepository()

	existingInstanceSettings := &domain.OrganizationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
			OrganizationScopedUsernames: gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.OrganizationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		OrganizationSettingsAttributes: domain.OrganizationSettingsAttributes{
			OrganizationScopedUsernames: gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeOrganization, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeOrganization, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeOrganization, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeOrganization, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeOrganization, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.NotificationSettingsRepository()

	settings := []*domain.NotificationSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.NotificationSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.NotificationSettingsRepository()

	settings := []*domain.NotificationSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
				PasswordChange: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.NotificationSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.NotificationSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.NotificationSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.NotificationSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.NotificationSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.NotificationSettingsRepository()

	existingSettings := &domain.NotificationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
			PasswordChange: gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.NotificationSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.NotificationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
					PasswordChange: gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.NotificationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
					PasswordChange: gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.NotificationSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
					PasswordChange: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.NotificationSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
					PasswordChange: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.NotificationSettingsRepository()

	existingInstanceSettings := &domain.NotificationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
			PasswordChange: gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.NotificationSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		NotificationSettingsAttributes: domain.NotificationSettingsAttributes{
			PasswordChange: gu.Ptr(true),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeNotification, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeNotification, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeNotification, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeNotification, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeNotification, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LegalAndSupportSettingsRepository()

	settings := []*domain.LegalAndSupportSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LegalAndSupportSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LegalAndSupportSettingsRepository()

	settings := []*domain.LegalAndSupportSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
				TOSLink:           gu.Ptr("https://host"),
				PrivacyPolicyLink: gu.Ptr("https://host"),
				HelpLink:          gu.Ptr("https://host"),
				SupportEmail:      gu.Ptr("email"),
				DocsLink:          gu.Ptr("https://host"),
				CustomLink:        gu.Ptr("https://host"),
				CustomLinkText:    gu.Ptr("linktext"),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LegalAndSupportSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.LegalAndSupportSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LegalAndSupportSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LegalAndSupportSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LegalAndSupportSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LegalAndSupportSettingsRepository()

	existingSettings := &domain.LegalAndSupportSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
			TOSLink:           gu.Ptr("https://host"),
			PrivacyPolicyLink: gu.Ptr("https://host"),
			HelpLink:          gu.Ptr("https://host"),
			SupportEmail:      gu.Ptr("email"),
			DocsLink:          gu.Ptr("https://host"),
			CustomLink:        gu.Ptr("https://host"),
			CustomLinkText:    gu.Ptr("linktext"),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LegalAndSupportSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LegalAndSupportSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
					TOSLink:           gu.Ptr("https://host"),
					PrivacyPolicyLink: gu.Ptr("https://host"),
					HelpLink:          gu.Ptr("https://host"),
					SupportEmail:      gu.Ptr("email"),
					DocsLink:          gu.Ptr("https://host"),
					CustomLink:        gu.Ptr("https://host"),
					CustomLinkText:    gu.Ptr("linktext"),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.LegalAndSupportSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
					TOSLink:           gu.Ptr("https://host"),
					PrivacyPolicyLink: gu.Ptr("https://host"),
					HelpLink:          gu.Ptr("https://host"),
					SupportEmail:      gu.Ptr("email"),
					DocsLink:          gu.Ptr("https://host"),
					CustomLink:        gu.Ptr("https://host"),
					CustomLinkText:    gu.Ptr("linktext"),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LegalAndSupportSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
					TOSLink:           gu.Ptr("https://host"),
					PrivacyPolicyLink: gu.Ptr("https://host"),
					HelpLink:          gu.Ptr("https://host"),
					SupportEmail:      gu.Ptr("email"),
					DocsLink:          gu.Ptr("https://host"),
					CustomLink:        gu.Ptr("https://host"),
					CustomLinkText:    gu.Ptr("linktext"),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LegalAndSupportSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
					TOSLink:           gu.Ptr("https://host"),
					PrivacyPolicyLink: gu.Ptr("https://host"),
					HelpLink:          gu.Ptr("https://host"),
					SupportEmail:      gu.Ptr("email"),
					DocsLink:          gu.Ptr("https://host"),
					CustomLink:        gu.Ptr("https://host"),
					CustomLinkText:    gu.Ptr("linktext"),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LegalAndSupportSettingsRepository()

	existingInstanceSettings := &domain.LegalAndSupportSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
			TOSLink:           gu.Ptr("https://host"),
			PrivacyPolicyLink: gu.Ptr("https://host"),
			HelpLink:          gu.Ptr("https://host"),
			SupportEmail:      gu.Ptr("email"),
			DocsLink:          gu.Ptr("https://host"),
			CustomLink:        gu.Ptr("https://host"),
			CustomLinkText:    gu.Ptr("linktext"),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LegalAndSupportSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LegalAndSupportSettingsAttributes: domain.LegalAndSupportSettingsAttributes{
			TOSLink:           gu.Ptr("https://host"),
			PrivacyPolicyLink: gu.Ptr("https://host"),
			HelpLink:          gu.Ptr("https://host"),
			SupportEmail:      gu.Ptr("email"),
			DocsLink:          gu.Ptr("https://host"),
			CustomLink:        gu.Ptr("https://host"),
			CustomLinkText:    gu.Ptr("linktext"),
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecretGeneratorSettingsRepository()

	settings := []*domain.SecretGeneratorSettings{
		createSecretGeneratorSettings(firstInstanceID, nil),
		createSecretGeneratorSettings(secondInstanceID, nil),
		createSecretGeneratorSettings(firstInstanceID, gu.Ptr(firstOrgID)),
		createSecretGeneratorSettings(secondInstanceID, gu.Ptr(secondOrgID)),
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.SecretGeneratorSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecretGeneratorSettingsRepository()

	settings := []*domain.SecretGeneratorSettings{
		createSecretGeneratorSettings(firstInstanceID, nil),
		createSecretGeneratorSettings(secondInstanceID, nil),
		createSecretGeneratorSettings(firstInstanceID, gu.Ptr(firstOrgID)),
		createSecretGeneratorSettings(secondInstanceID, gu.Ptr(secondOrgID)),
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.SecretGeneratorSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.SecretGeneratorSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.SecretGeneratorSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.SecretGeneratorSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.SecretGeneratorSettings{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecretGeneratorSettingsRepository()

	existingSettings := createSecretGeneratorSettings(instanceID, gu.Ptr(orgID))

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.SecretGeneratorSettings
		wantErr  error
	}{
		{
			name:     "create instance",
			settings: createSecretGeneratorSettings(instanceID, nil),
		},
		{
			name:     "update organization",
			settings: createSecretGeneratorSettings(instanceID, gu.Ptr(orgID)),
		},
		{
			name:     "non-existing instance",
			settings: createSecretGeneratorSettings("foo", nil),
			wantErr:  new(database.ForeignKeyError),
		},
		{
			name:     "non-existing org",
			settings: createSecretGeneratorSettings(instanceID, gu.Ptr("foo")),
			wantErr:  new(database.ForeignKeyError),
		},
		{
			name:     "nil settings",
			settings: nil,
			wantErr:  database.ErrInvalidChangeTarget,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecretGeneratorSettingsRepository()

	existingInstanceSettings := createSecretGeneratorSettings(instanceID, nil)
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := createSecretGeneratorSettings(instanceID, gu.Ptr(orgID))
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func createDefaultSecretGeneratorAttrs() domain.SecretGeneratorAttrs {
	return domain.SecretGeneratorAttrs{
		Length:              gu.Ptr(uint(32)),
		IncludeUpperLetters: gu.Ptr(true),
		IncludeLowerLetters: gu.Ptr(true),
		IncludeDigits:       gu.Ptr(true),
		IncludeSymbols:      gu.Ptr(false),
	}
}

func createDefaultSecretGeneratorAttrsWithExpiry(expiry time.Duration) domain.SecretGeneratorAttrsWithExpiry {
	return domain.SecretGeneratorAttrsWithExpiry{
		Expiry: gu.Ptr(expiry),
		SecretGeneratorAttrs: domain.SecretGeneratorAttrs{
			Length:              gu.Ptr(uint(32)),
			IncludeUpperLetters: gu.Ptr(true),
			IncludeLowerLetters: gu.Ptr(true),
			IncludeDigits:       gu.Ptr(true),
			IncludeSymbols:      gu.Ptr(false),
		},
	}
}

func createDefaultSecretGeneratorSettingsAttributes() domain.SecretGeneratorSettingsAttributes {
	attrs := createDefaultSecretGeneratorAttrs()
	oneDay := 24 * time.Hour
	thirtyMinutes := 30 * time.Minute
	tenMinutes := 10 * time.Minute

	return domain.SecretGeneratorSettingsAttributes{
		ClientSecret: &domain.ClientSecretAttributes{
			SecretGeneratorAttrs: attrs,
		},
		InitializeUserCode: &domain.InitializeUserCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(oneDay),
		},
		EmailVerificationCode: &domain.EmailVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PhoneVerificationCode: &domain.PhoneVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PasswordVerificationCode: &domain.PasswordVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PasswordlessInitCode: &domain.PasswordlessInitCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(tenMinutes),
		},
		DomainVerification: &domain.DomainVerificationAttributes{
			SecretGeneratorAttrs: attrs,
		},
		OTPSMS: &domain.OTPSMSAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		OTPEmail: &domain.OTPEmailAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
	}
}

func createSecretGeneratorSettings(instanceID string, organizationID *string) *domain.SecretGeneratorSettings {
	return &domain.SecretGeneratorSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: organizationID,
		},
		SecretGeneratorSettingsAttributes: createDefaultSecretGeneratorSettingsAttributes(),
	}
}
